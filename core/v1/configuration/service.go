package configuration

import (
	"context"
	"encoding/json"
	http "github.com/PicPay/ms-data-crawler/core/v1"
	"github.com/PicPay/ms-data-crawler/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"sync"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../configuration/mock_gateway.go -package=configuration -self_package=github.com/PicPay/ms-data-crawler/core/v1/configuration . Gateway
type Gateway interface {
	Find(ctx context.Context, in interface{}) (*Configuration, error)
}

type ServiceChannel struct {
	data interface{}
	name string
}

type Service struct {
	Gateway
	*Configuration
	http.HttpClient
}

func NewService(gateway Gateway, httpClient http.HttpClient) *Service {
	return &Service{gateway, &Configuration{}, httpClient}
}

func (s *Service) Fetch(ctx context.Context, identifier string, headers map[string]string) (*AssembledScreen, error) {
	log.Info("Started", &log.LogContext{
		"Function":   "Fetch",
		"identifier": identifier,
		"headers":    headers,
	})

	data := make(chan ServiceChannel)
	var dataServices []interface{}

	configuration, err := s.Gateway.Find(ctx, bson.M{"identifier": identifier})
	if err != nil {
		return nil, err
	}

	s.Configuration = configuration

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

	assembledScreen := setScreenTemplate(
		identifier,
		configuration.CreatedAt,
		configuration.UpdatedAt,
		dataServices,
	)

	log.Info("Finished", &log.LogContext{
		"Function": "Fetch",
	})

	return &assembledScreen, err
}

func (s *Service) getPreData(ctx context.Context, source ServiceRequest, headers map[string]string) (midgardRows []map[string]string, err error) {
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
		"function": "getPreData",
		"url":      url,
	})

	resp, err := s.HttpClient.Send(ctx, httpConfiguration, nil)
	if err != nil || resp.Status != 200 {
		log.Error("Error getting data from:", err, &log.LogContext{
			"function": "getPreData",
			"url":      url,
			"status":   resp.Status,
			"error":    err,
		})
		return
	}

	log.Debug("Finished data from:", &log.LogContext{
		"function": "getPreData",
		"url":      source.Url,
		"status":   resp.Status,
		"resp":     string(resp.Body),
	})

	if err = json.Unmarshal(resp.Body, &midgardRows); err == nil {
		return
	}

	log.Warn("Error from Unmasharling data from - needs []map[string] or responseFormat json", &log.LogContext{
		"url":   source.Url,
		"error": err,
	})

	if s.Source.HasResponseFormat() {
		midgardRows, err = s.formatResponseToMarkupJson(resp.Body, s.Source.ResponseFormat)
	}

	return
}

func (s *Service) formatResponseToMarkupJson(body []byte, markupJson string) (response []map[string]string, err error) {
	var serviceResponse interface{}
	var newBody []byte

	if err := json.Unmarshal(body, &serviceResponse); err != nil {
		log.Error("Error from Unmasharling data to interface{} from:", err, &log.LogContext{
			"Function": "formatResponseToMarkupJson",
			"error":    err,
		})

		return response, err
	}

	tpl, err := NewTemplateParser("newJson", markupJson)

	if err != nil {
		log.Error("Invalid template", err, &log.LogContext{
			"error":    err,
			"template": markupJson,
		})

		return
	}

	newBody, err = tpl.Parse(serviceResponse)

	if err != nil {
		log.Error("Error from parsing template into data:", err, &log.LogContext{
			"Class": "formatResponseToMarkupJson",
			"error": err,
			"data":  string(newBody),
		})

		return response, err
	}

	if err := json.Unmarshal(newBody, &response); err != nil {
		log.Error("Error from Unmasharling newBody data - Check template for typo's:", err, &log.LogContext{
			"Class": "formatResponseToMarkupJson",
			"error": err,
		})

		return response, err
	}

	return
}

func (s *Service) getDataFromService(
	ctx context.Context,
	data chan ServiceChannel,
	midgardRows []map[string]string,
	Crawlers []ServiceRequest,
	headers map[string]string) {

	log.Info("Started", &log.LogContext{
		"function":      "getDataFromService",
		"Service count": len(midgardRows),
	})

	var wg sync.WaitGroup
	wg.Add(len(midgardRows))

	for _, midgardData := range midgardRows {
		httpClient := s.HttpClient
		go func(wg *sync.WaitGroup, data chan ServiceChannel, Crawlers []ServiceRequest, midgardData map[string]string, headers map[string]string, httpClient http.HttpClient) {
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

			httpConfiguration := http.HttpConfiguration{
				URL:    url,
				Method: "GET",
			}

			resp, err := httpClient.Send(ctx, httpConfiguration, nil)
			if err != nil || resp.Status != 200 {
				log.Error("Error getting data from:", err, &log.LogContext{
					"Class":  "getDataFromService",
					"url":    url,
					"status": resp.Status,
					"error":  err,
				})
				return
			}

			body := resp.Body
			log.Debug("Finished data from:", &log.LogContext{
				"Class":  "getDataFromService",
				"url":    url,
				"status": resp.Status,
				"resp":   string(body),
			})

			var result interface{}
			if err = json.Unmarshal(body, &result); err != nil {
				log.Error("Error from Unmasharling data from:", err, &log.LogContext{
					"Class": "getDataFromService",
					"url":   url,
					"error": err,
				})

				return
			}

			if len(body) > 0 {
				data <- ServiceChannel{result, time.Now().String()}
			}
		}(&wg, data, Crawlers, midgardData, headers, httpClient)
	}

	wg.Wait()
	close(data)

	log.Info("Finished", &log.LogContext{
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
