package configuration

import (
	"github.com/PicPay/ms-data-crawler/pkg/server"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (h *Handler) Load(r *gin.RouterGroup, server *server.Server) error {
	dataRepository := NewRepository(server.DB)

	assemblerService := NewService(dataRepository)
	controller := NewController(assemblerService)

	r.GET("/configuration/:Identifier", controller.Fetch)

	return nil
}
