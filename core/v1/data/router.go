package data

import (
	"github.com/PicPay/ms-data-formatter/pkg/server"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (h *Handler) Load(r *gin.RouterGroup, server *server.Server) error {
	dataRepository := NewRepository(server.DB)

	assemblerService := NewService(dataRepository)
	controller := NewController(assemblerService)

	r.GET("/data/:ConsumerId", controller.Format)

	return nil
}
