package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/PicPay/ms-data-crawler/core/v1"
	"github.com/PicPay/ms-data-crawler/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"sync"
	"time"
)

type Gateway interface {
	Find(ctx context.Context, in interface{}) (*Data, error)
}

type ServiceChannel struct {
	data interface{}
	name string
}

type Service struct {
	Gateway
}

func NewService(gateway Gateway) *Service {
	return &Service{gateway}
}

func (s *Service) Fetch(ctx context.Context, identifier string, headers map[string]string) (*AssembledScreen, error) {
	log.Info("Started", &log.LogContext{
		"Function":   "Format",
		"identifier": identifier,
		"headers":    headers,
	})

	data := make(chan ServiceChannel)
	var dataServices []interface{}

	configuration, err := s.Gateway.Find(ctx, bson.M{"identifier": identifier})
	fmt.Println(configuration.Source)
	if err != nil {
		return nil, err
	}

	midgardRows, err := s.getPreData(ctx, configuration.Source, headers)
	if err != nil {
		return nil, err
	}

	go s.getDataFromService(ctx, data, midgardRows, configuration.Crawler, headers)

	for service := range data {
		if arr, ok := service.data.([]interface{}); ok {
			dataServices = append(dataServices, arr)
		} else {
			dataServices = append(dataServices, service.data)
		}
	}

	fmt.Println("dataServices", dataServices)

	assembledScreen := setScreenTemplate(
		identifier,
		configuration.CreatedAt,
		configuration.UpdatedAt,
		dataServices,
	)

	log.Info("Finished", &log.LogContext{
		"Function": "Format",
	})

	return &assembledScreen, err
}

func (s *Service) getPreData(ctx context.Context, source ServiceRequest, headers map[string]string) (midgardRows []map[string]string, err error) {

	httpClient := http.NewHttpAdapterWithOptions(10 * time.Second)

	url := source.Url

	// if service mapping matches with any header key,
	// we replace with its value
	if source.HasMapping() {
		url = replaceUrlMarkup(url, source.Mapping, headers)
	}

	httpConfiguration := http.HttpConfiguration{
		URL: url,
	}

	log.Info("Getting data from:", &log.LogContext{
		"Class": "AssemblerService",
		"url":   url,
	})

	resp, err := httpClient.Send(ctx, httpConfiguration, nil)
	if err != nil || resp.Status != 200 {
		log.Error("Error getting data from:", err, &log.LogContext{
			"Class":  "AssemblerService",
			"url":    url,
			"status": resp.Status,
			"error":  err,
		})
		return
	}

	log.Debug("Finished data from:", &log.LogContext{
		"Class":  "AssemblerService",
		"url":    source.Url,
		"status": resp.Status,
		"resp":   string(resp.Body),
	})

	if err = json.Unmarshal(resp.Body, &midgardRows); err != nil {
		log.Error("Error from Unmasharling data from:", err, &log.LogContext{
			"Class": "AssemblerService",
			"url":   source.Url,
			"error": err,
		})

		return
	}
	/*
		fmt.Println("midgardRows", midgardRows)
		var ttest []map[string]string
		for key, value := range jsonMap {
			if reflect.TypeOf(value).String() == "map[string]interface {}" {
				tt1 := make(map[string]string)
				for key1, value1 := range value.(map[string]interface{}) {
					tt1[key1] = value1.(string)

				}
				ttest = append(ttest, tt1)
			}
		}

		fmt.Println("Eta porra", ttest, ttest[0]["id"])
	*/

	return
}

func (s *Service) getDataFromService(
	ctx context.Context,
	data chan ServiceChannel,
	midgardRows []map[string]string,
	Crawlers []ServiceRequest,
	headers map[string]string) {

	log.Info("Started", &log.LogContext{
		"Class":         "AssemblerService",
		"function":      "getDataFromService",
		"Service count": len(midgardRows),
	})

	var wg sync.WaitGroup
	wg.Add(len(midgardRows))

	for _, midgardData := range midgardRows {
		fmt.Println("midgardData", midgardData)
		go func(wg *sync.WaitGroup, data chan ServiceChannel, Crawlers []ServiceRequest, midgardData map[string]string, headers map[string]string) {
			defer wg.Done()

			crawler := validateCrawler(Crawlers, midgardData["type"])

			if crawler.Url == "" {
				return
			}

			url := crawler.Url

			// if service mapping matches with any header key,
			// we replace with its value
			if crawler.HasMapping() {
				url = replaceUrlMarkup(url, crawler.Mapping, headers)
				url = replaceUrlMarkup(url, crawler.Mapping, midgardData)
			}

			fmt.Println("Crawler", url)

			log.Debug("Request header", &log.LogContext{
				"Class": "AssemblerService",
				"url":   url,
			})

			httpClient := http.NewHttpAdapterWithOptions(10 * time.Second)
			httpConfiguration := http.HttpConfiguration{
				URL:    url,
				Method: "GET",
			}

			resp, err := httpClient.Send(ctx, httpConfiguration, nil)
			if err != nil || resp.Status != 200 {
				log.Error("Error getting data from:", err, &log.LogContext{
					"Class":  "AssemblerService",
					"url":    url,
					"status": resp.Status,
					"error":  err,
				})
				return
			}

			body := resp.Body
			log.Debug("Finished data from:", &log.LogContext{
				"Class":  "AssemblerService",
				"url":    url,
				"status": resp.Status,
				"resp":   string(body),
			})

			log.Debug("Applying body transformations:", &log.LogContext{})

			var result interface{}
			if err = json.Unmarshal(body, &result); err != nil {
				log.Error("Error from Unmasharling data from:", err, &log.LogContext{
					"Class": "AssemblerService",
					"url":   url,
					"error": err,
				})

				return
			}

			if len(body) > 0 {
				data <- ServiceChannel{result, time.Now().String()}
			}
		}(&wg, data, Crawlers, midgardData, headers)
	}

	wg.Wait()
	close(data)

	log.Info("Finished", &log.LogContext{
		"Class":    "AssemblerService",
		"function": "getDataFromService",
	})
}

func validateCrawler(crawlers []ServiceRequest, crawlerType string) ServiceRequest {
	for _, crawler := range crawlers {
		if crawler.Validation[0].Value == crawlerType {
			return crawler
		}
	}
	return ServiceRequest{}
}

func replaceUrlMarkup(url string, mapping []KeyValue, keys map[string]string) (newUrl string) {
	newUrl = url
	for _, markup := range mapping {
		fmt.Println("markup", url, markup.Index, markup.Value, keys[markup.Value])
		if keys[markup.Value] != "" {
			newUrl = strings.ReplaceAll(url, markup.Index, keys[markup.Value])
			return
		}
	}

	return
}

func setScreenTemplate(identifier string, createdAt, updatedAt time.Time, body interface{}) AssembledScreen {
	var assembledScreen AssembledScreen

	assembledScreen.Identifier = identifier
	assembledScreen.CreatedAt = createdAt
	assembledScreen.UpdatedAt = updatedAt
	assembledScreen.Data = body

	return assembledScreen
}
