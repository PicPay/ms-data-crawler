package assembler

import (
	"errors"
	"github.com/PicPay/picpay-dev-ms-template-manager/core/v2/template"
	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
	"strconv"
	"strings"

	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Controller struct {
	*Service
}

type Params struct {
	ExternalData []KeyValue `json:"external_data"`
}

// headers that can be forwarded to search service
var proxyHeaders = []string{
	"App_Version",
	"Device_Os",
	"Latitude",
	"Location_Timestamp",
	"Longitude",
	"Timezone",
	"Token",
	"Consumer_Id",
	"Area_Code",
	"X-Request-ID",
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

func (c *Controller) getRequestParameters(ctx *gin.Context) (reqParams Params) {
	if strings.ToLower(ctx.Request.Method) == "post" {
		ctx.ShouldBindBodyWith(&reqParams, binding.JSON)
	}
	return
}

func (c *Controller) AssembleScreen(ctx *gin.Context) {
	log.WithContext(ctx)

	screenName := ctx.Param("ScreenName")
	if screenName == "" {
		http.BadRequest(ctx, errors.New("screen name is required"))
		return
	}

	appVersionSplit := strings.Split(ctx.Request.Header.Get("app_version"), ".")
	deviceOs := ctx.Request.Header.Get("device_os")
	if deviceOs == "" || len(appVersionSplit) < 3 {
		http.BadRequest(ctx, errors.New("app_version and device_os is required"))
		return
	}

	major, _ := strconv.ParseUint(appVersionSplit[0], 10, 64)
	minor, _ := strconv.ParseUint(appVersionSplit[1], 10, 64)
	patch, _ := strconv.ParseUint(appVersionSplit[2], 10, 64)

	listRequest := template.ListRequest{
		DeviceOs: ctx.Request.Header.Get("device_os"),
		AppVersion: template.SemanticVersion{
			Major: major,
			Minor: minor,
			Patch: patch,
		},
		ScreenName: screenName,
	}

	headers := c.getProxyHeaders(ctx)
	params := c.getRequestParameters(ctx)

	screen, err := c.Service.AssemblePage(ctx, listRequest, headers, params.ExternalData)
	if err != nil {
		http.InternalServerError(ctx, err)
		return
	}

	http.Ok(ctx, screen)
}
