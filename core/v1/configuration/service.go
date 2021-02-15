package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/PicPay/ms-data-crawler/core/v1"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"sync"
	"time"
	/*"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	*/
	"github.com/PicPay/ms-data-crawler/pkg/log"
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

func (s *Service) Format(ctx context.Context, identifier string, headers map[string]string) (*AssembledScreen, error) {
	log.Info("Started", &log.LogContext{
		"Function":   "Format",
		"identifier": identifier,
		"headers":    headers,
	})

	data := make(chan ServiceChannel)
	var dataServices []interface{}

	crawler, err := s.Gateway.Find(ctx, bson.M{"identifier": identifier})
	fmt.Println(data)
	if err != nil {
		return nil, err
	}

	midgardRows, err := s.getDataFromDatalake(ctx, crawler.Url)
	if err != nil {
		return nil, err
	}

	fmt.Println("midgardData", midgardRows)

	go s.getServices(ctx, data, midgardRows, crawler.UrlSource)

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
		crawler.CreatedAt,
		crawler.UpdatedAt,
		dataServices,
	)

	log.Info("Finished", &log.LogContext{
		"Function": "Format",
	})

	return &assembledScreen, err
}

func (s *Service) getDataFromDatalake(ctx context.Context, url string) (midgardRows []Midgard, err error) {

	httpClient := http.NewHttpAdapterWithOptions(10 * time.Second)
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
		"url":    url,
		"status": resp.Status,
		"resp":   string(resp.Body),
	})

	log.Debug("Applied body transformations:", &log.LogContext{})

	if err = json.Unmarshal(resp.Body, &midgardRows); err != nil {
		log.Error("Error from Unmasharling data from:", err, &log.LogContext{
			"Class": "AssemblerService",
			"url":   url,
			"error": err,
		})
		return
	}

	return
}

func (s *Service) getServices(ctx context.Context, data chan ServiceChannel, midgardRows []Midgard, dgUrl string) {
	log.Info("Started", &log.LogContext{
		"Class":         "AssemblerService",
		"function":      "getServices",
		"Service count": len(midgardRows),
	})

	var wg sync.WaitGroup
	wg.Add(len(midgardRows))

	for _, midgardData := range midgardRows {
		fmt.Println("midgardData", midgardData.getId())
		dgId := midgardData.getId()
		go func(wg *sync.WaitGroup, data chan ServiceChannel, dgUrl, dgId string) {
			defer wg.Done()

			url := dgUrl

			log.Debug("Request header", &log.LogContext{
				"Class": "AssemblerService",
				"url":   url,
			})

			url = strings.ReplaceAll(url, ":dgId", dgId)

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
		}(&wg, data, dgUrl, dgId)
	}

	wg.Wait()
	close(data)

	log.Info("Finished", &log.LogContext{
		"Class":    "AssemblerService",
		"function": "getServices",
	})
}

func setScreenTemplate(identifier string, createdAt, updatedAt time.Time, body interface{}) AssembledScreen {
	var assembledScreen AssembledScreen

	assembledScreen.Identifier = identifier
	assembledScreen.CreatedAt = createdAt
	assembledScreen.UpdatedAt = updatedAt
	assembledScreen.Data = body

	return assembledScreen
}
