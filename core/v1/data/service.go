package data

import (
	"context"
	/*"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	http "github.com/PicPay/ms-data-formatter/core/v1"*/
	"github.com/PicPay/ms-data-formatter/pkg/log"
)

type Gateway interface {
	Find(ctx context.Context, in interface{}) (*Data, error)
}

type Service struct {
	Gateway
}

func NewService(gateway Gateway) *Service {
	return &Service{gateway}
}

func (s *Service) Format(ctx context.Context, consumerId string) (*AssembledScreen, error) {
	log.Info("Started", &log.LogContext{
		"Function": "Format",
	})

	/*
		go s.getServices(ctx, data, template.Services, headers, params)

		rawTemplate, err := s.getComponents(ctx, template, params)
		if err != nil {
			return nil, err
		}

		for service := range data {
			if arr, ok := service.data.([]interface{}); ok {
				dataServices[service.name] = arr
			} else {
				dataServices[service.name] = service.data
			}
		}

		cleanedJSON, err := parseTemplateData(rawTemplate, dataServices)
		if err != nil {
			return nil, err
		}

		assembledScreen := setScreenTemplate(
			template.CreatedAt,
			template.UpdatedAt,
			cleanedJSON,
		)

		log.Info("Finished", &log.LogContext{
			"Function":   "Format",
		})

		return &assembledScreen, err
	*/
	return nil, nil
}

/*
func (s *Service) getServices(ctx context.Context, data chan ServiceChannel, services []template.ServiceRequest, headers map[string]string, params []KeyValue) {
	log.Info("Started", &log.LogContext{
		"Class":         "AssemblerService",
		"function":      "getServices",
		"Service count": len(services),
	})

	var wg sync.WaitGroup
	wg.Add(len(services))

	for _, service := range services {
		go func(wg *sync.WaitGroup, data chan ServiceChannel, service template.ServiceRequest, httpHeaders map[string]string, params []KeyValue) {
			defer wg.Done()

			headers := make(map[string]string)

			// don't forward http headers if the service is external
			if !service.External {
				headers = httpHeaders
			}

			if len(service.Headers) > 0 {
				headers = mergeHeadersMap(headers, service.Headers)
			}

			url := service.Url
			// if service mapping matches with any header key,
			// we replace with its value
			if service.HasMapping() {
				url = replaceUrlMarkup(url, service.Mapping, headers)
			}

			// if a service url supports param matching like: `www.x.com/:id`,
			// we search for the key "id" in ExternalParams,
			// if any found we replace with its value
			for _, param := range params {
				if key := ":" + param.Key; strings.Contains(url, key) {
					url = strings.ReplaceAll(url, key, param.Value)
				}
			}

			log.Debug("Request header", &log.LogContext{
				"Class":   "AssemblerService",
				"url":     url,
				"headers": headers,
				"service": service,
			})

			httpClient := http.NewHttpAdapterWithOptions(10 * time.Second)
			httpConfiguration := http.HttpConfiguration{
				URL:         url,
				Headers:     headers,
				Method:      service.Method,
				ContentType: service.ContentType,
			}

			if httpConfiguration.Method == "" {
				httpConfiguration.Method = "GET"
			}

			log.Info("Getting data from:", &log.LogContext{
				"Class":   "AssemblerService",
				"url":     url,
				"headers": headers,
				"Method":  service.Method,
			})

			resp, err := httpClient.Send(ctx, httpConfiguration, nil)
			if err != nil || resp.Status != statusOK {
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

			for _, onRes := range service.OnResponse {
				for _, transform := range onRes.BodyTransform {
					transformFn, ok := ResponseTransformer[transform.Name]
					if !ok {
						log.Debug(fmt.Sprintf("body transformations not found: %s", transform.Name), &log.LogContext{})
						continue
					}
					body = transformFn(body, transform.Args)
					log.Debug(fmt.Sprintf("body transformation applied: %s", transform.Name), &log.LogContext{})
				}
			}

			log.Debug("Applied body transformations:", &log.LogContext{})

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
				data <- ServiceChannel{result, service.Name}
			}
		}(&wg, data, service, headers, params)
	}

	wg.Wait()
	close(data)

	log.Info("Finished", &log.LogContext{
		"Class":    "AssemblerService",
		"function": "getServices",
	})
}


func setScreenTemplate(identifier, version string, createdAt, updatedAt time.Time, body interface{}) AssembledScreen {
	var assembledScreen AssembledScreen

	newVersion, err := strconv.ParseFloat(version, 2)
	if err != nil {
		log.Error("Error on casting version to float", err, &log.LogContext{"Class": "AssemblerService"})
	}

	assembledScreen.Identifier = identifier
	assembledScreen.Version = int(newVersion)
	assembledScreen.CreatedAt = createdAt
	assembledScreen.UpdatedAt = updatedAt
	assembledScreen.Body = body

	return assembledScreen
}
*/
