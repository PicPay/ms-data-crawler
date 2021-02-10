package assembler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	http "github.com/PicPay/picpay-dev-ms-template-manager/core/v2"
	"github.com/PicPay/picpay-dev-ms-template-manager/core/v2/component"
	"github.com/PicPay/picpay-dev-ms-template-manager/core/v2/template"
	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
)

type Service struct {
	templateService  *template.Service
	componentService *component.Service
}

type ServiceChannel struct {
	data interface{}
	name string
}

const (
	maxTemplateParams = 100
	statusOK          = 200
	jsonTemplate      = `{
	"header": %s,
	"components": [%s]
}`
)

func NewService(templateService *template.Service, componentService *component.Service) *Service {
	return &Service{
		templateService:  templateService,
		componentService: componentService,
	}
}

func (s *Service) AssemblePage(ctx context.Context, request template.ListRequest, headers map[string]string, params []KeyValue) (*AssembledScreen, error) {
	log.Info("Started", &log.LogContext{
		"Class":      "AssemblerService",
		"Function":   "AssemblePage",
		"screenName": request.ScreenName,
	})

	data := make(chan ServiceChannel)
	dataServices := make(map[string]interface{})

	template, err := s.getTemplate(ctx, request)
	if err != nil {
		return nil, err
	}

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
		template.ScreenName,
		template.Version,
		template.CreatedAt,
		template.UpdatedAt,
		cleanedJSON,
	)

	log.Info("Finished", &log.LogContext{
		"Class":      "AssemblerService",
		"Function":   "AssemblePage",
		"screenName": request.ScreenName,
	})

	return &assembledScreen, err
}

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

func replaceUrlMarkup(url string, mapping []template.KeyValue, headers map[string]string) (newUrl string) {
	for _, markup := range mapping {
		newUrl = strings.ReplaceAll(url, markup.Index, headers[markup.Value])
	}

	return
}

func (s *Service) getComponents(ctx context.Context, templateRow *template.Template, params []KeyValue) (string, error) {
	log.Info("Started", &log.LogContext{
		"Class":            "AssemblerService",
		"Function":         "getComponents",
		"components count": len(templateRow.Components),
	})

	log.Debug("Searching header", &log.LogContext{
		"Class":      "AssemblerService",
		"identifier": templateRow.Header.ComponentId,
	})

	headerComponent, err := s.componentService.FindBy(ctx, "identifier", templateRow.Header.ComponentId)
	if err != nil {
		log.Error("Couldnt find header", err, &log.LogContext{
			"Class":                "AssemblerService",
			"component identifier": templateRow.Header.ComponentId,
		})
		return "", err
	}

	header := headerComponent.Markup
	for _, keyValue := range templateRow.Header.Mapping {
		header = strings.ReplaceAll(header, keyValue.Index, keyValue.Value)
	}

	var componentsCollection []string
	for _, componentRow := range templateRow.Components {
		if componentRow.IsAssembled() {
			componentsCollection = append(componentsCollection, componentRow.Assembled)
			continue
		}

		component, err := s.componentService.FindBy(ctx, "identifier", componentRow.ComponentId)
		if err != nil {
			log.Error("Couldnt find component", err, &log.LogContext{
				"Class":      "AssemblerService",
				"identifier": componentRow.ComponentId,
			})

			return "", err
		}

		componentMarkup := component.Markup
		for _, mapping := range componentRow.Mapping {
			// parameters mappings are prefixed with a "#" and its values in the collection are represented as default
			if strings.HasPrefix(mapping.Index, "#") {
				if len(params) > 0 && len(params) < maxTemplateParams {
					key := mapping.Index[1:]
					for _, param := range params {
						if param.Key == key {
							componentMarkup = strings.ReplaceAll(componentMarkup, mapping.Index, param.Value)
							break
						}
					}
				}
			}

			componentMarkup = strings.ReplaceAll(componentMarkup, mapping.Index, mapping.Value)
		}

		componentMarkup = strings.ReplaceAll(componentMarkup, "@@service", componentRow.Service)

		componentsCollection = append(componentsCollection, componentMarkup)
	}

	components := strings.Join(componentsCollection, ",")
	markup := fmt.Sprintf(jsonTemplate, header, components)

	log.Info("Finished", &log.LogContext{"Class": "AssemblerService", "Function": "getComponents"})

	return markup, nil
}

func (s *Service) getTemplate(ctx context.Context, request template.ListRequest) (*template.Template, error) {
	log.Debug("Started ", &log.LogContext{
		"Class":       "AssemblerService",
		"getTemplate": "getTemplate",
		"screenName":  request.ScreenName,
	})

	template, err := s.templateService.FindBy(ctx, request)
	if err != nil {
		log.Error("Could not found Template", err, &log.LogContext{
			"Class":      "AssemblerService",
			"screenName": request.ScreenName,
		})
		return nil, err
	}

	log.Debug("Template found!", &log.LogContext{"Class": "AssemblerService", "screenName": request.ScreenName})

	if len(template.Components) == 0 {
		log.Error("Template with no components", err, &log.LogContext{
			"Class":      "AssemblerService",
			"screenName": request.ScreenName,
		})
		return nil, err
	}

	return template, nil
}

func parseTemplateData(template string, data map[string]interface{}) (interface{}, error) {
	log.Info("Started", &log.LogContext{
		"Class":    "AssemblerService",
		"function": "parseTemplateData",
	})

	tpl, err := NewTemplateParser("screen", template)
	if err != nil {
		return nil, err
	}

	body, err := tpl.Parse(data)

	if err != nil {
		log.Error("Couldnt set data on template, file outputted", err, &log.LogContext{
			"Class":        "AssemblerService",
			"file_content": string(body),
		})
		return nil, err
	}

	return cleanTemplate(string(body))
}

func mergeHeadersMap(headers1 map[string]string, headers2 []template.KeyValue) map[string]string {
	newHeader := make(map[string]string)

	for key, value := range headers1 {
		newHeader[key] = value
	}

	for _, header := range headers2 {
		if header.Value != "" {
			newHeader[header.Index] = header.Value
		}
	}

	return newHeader
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
