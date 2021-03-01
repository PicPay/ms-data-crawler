package configuration

import (
	httpAdapter "github.com/PicPay/ms-data-crawler/core/v1"
	"github.com/PicPay/ms-data-crawler/pkg/server"
	"github.com/gin-gonic/gin"
	"time"
)

type Handler struct{}

func (h *Handler) Load(r *gin.RouterGroup, server *server.Server) error {
	dataRepository := NewRepository(server.DB)
	controller := NewController(
		NewService(
			dataRepository,
			httpAdapter.NewHttpAdapterWithOptions(1*time.Minute),
		),
	)

	r.GET("/configuration/:Identifier", controller.Fetch)

	return nil
}
