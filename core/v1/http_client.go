package v1

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/PicPay/ms-data-formatter/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type HttpAdapter struct {
	cli *resty.Client
}

func NewDefaultHttpAdapter() HttpAdapter {
	return HttpAdapter{
		cli: resty.New().SetTimeout(3 * time.Second),
	}
}

func NewHttpAdapterWithOptions(timeout time.Duration) HttpAdapter {
	return HttpAdapter{
		cli: resty.New().SetTimeout(timeout),
	}
}

func (adapter HttpAdapter) Send(ctx context.Context, configuration HttpConfiguration, request interface{}) (response HttpResponse, err error) {

	configURL := configuration.URL

	if _, err = url.ParseRequestURI(configURL); err != nil {
		err = errors.New(fmt.Sprintf("configuration url is not a valid url: %s", configURL))
		return
	}

	req := adapter.cli.SetHeaders(configuration.Headers).SetQueryParams(configuration.QueryParams)

	log.Info("Sending Http Request", &log.LogContext{
		"Method":       configuration.Method,
		"Headers":      configuration.Headers,
		"Content-Type": configuration.ContentType,
		"URL":          configuration.URL,
		"QueryParams":  configuration.QueryParams,
		"Body":         configuration.Data,
	})

	var resp *resty.Response
	resp, err = req.R().SetBody(request).ForceContentType(configuration.ContentType).SetContext(ctx).Execute(configuration.Method, configURL)

	txn := newrelic.FromContext(ctx)
	if txn != nil {
		txn = txn.NewGoroutine()
		segment := &newrelic.ExternalSegment{
			StartTime: txn.StartSegmentNow(),
			URL:       configuration.URL,
			Library:   "resty",
		}
		defer segment.End()
	}

	defer adapter.cli.SetCloseConnection(true)

	if resp != nil {
		response.Status = resp.StatusCode()
		response.Body = resp.Body()
	}

	return response, err
}
