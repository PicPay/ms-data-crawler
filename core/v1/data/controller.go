package data

import (
	"errors"
	"fmt"
	"github.com/PicPay/ms-data-formatter/pkg/http"
	"github.com/PicPay/ms-data-formatter/pkg/log"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	*Service
}

func NewController(service *Service) *Controller {
	return &Controller{service}
}

func (c *Controller) Format(ctx *gin.Context) {
	log.WithContext(ctx)

	consumerId := ctx.Param("ConsumerId")

	fmt.Println("ConsumerId", consumerId)
	if consumerId == "" {
		http.BadRequest(ctx, errors.New("Consumer Id is required"))
		return
	}

	jsonData, err := c.Service.Format(ctx, consumerId)
	if err != nil {
		http.InternalServerError(ctx, err)
		return
	}

	http.Ok(ctx, jsonData)
}
