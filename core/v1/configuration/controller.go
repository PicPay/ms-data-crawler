package configuration

import (
	"errors"
	"github.com/PicPay/ms-data-crawler/pkg/http"
	"github.com/PicPay/ms-data-crawler/pkg/log"
	"github.com/gin-gonic/gin"
)

// headers that can be forwarded to search service
var proxyHeaders = []string{
	"Consumer_Id",
}

type Controller struct {
	*Service
}

func NewController(service *Service) *Controller {
	return &Controller{service}
}

func (c *Controller) getProxyHeaders(ctx *gin.Context) map[string]string {
	headers := make(map[string]string)
	reqHeaders := ctx.Request.Header

	for _, proxyHeader := range proxyHeaders {
		value := reqHeaders.Get(proxyHeader)
		if value == "" {
			continue
		}

		headers[proxyHeader] = value
	}

	return headers
}

func (c *Controller) Format(ctx *gin.Context) {
	log.WithContext(ctx)

	identifier := ctx.Param("Identifier")

	if identifier == "" {
		http.BadRequest(ctx, errors.New("Identifier is required"))
		return
	}

	headers := c.getProxyHeaders(ctx)

	jsonData, err := c.Service.Format(ctx, identifier, headers)
	if err != nil {
		http.InternalServerError(ctx, err)
		return
	}

	http.Ok(ctx, jsonData)
}
