package assembler

import (
	"github.com/PicPay/picpay-dev-ms-template-manager/core/v2/component"
	"github.com/PicPay/picpay-dev-ms-template-manager/core/v2/template"
	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/server"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func (h *Handler) Load(r *gin.RouterGroup, server *server.Server) error {
	templateRepository := template.NewRepository(server.DB)
	componentRepository := component.NewRepository(server.DB)
	componentService := component.NewService(componentRepository)

	assemblerService := NewService(
		template.NewService(templateRepository, componentService),
		componentService,
	)
	controller := NewController(assemblerService)

	r.GET("/assembler/:ScreenName", controller.AssembleScreen)
	r.POST("/assembler/:ScreenName", controller.AssembleScreen)

	return nil
}
