// Package health is a general purpose health check http middleware
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Alive bool `json:"alive"`
}

func Handler(path string) func(w http.ResponseWriter, r *http.Request) bool {
	return func(w http.ResponseWriter, r *http.Request) bool {
		return r.URL.Path == path
	}
}

func GinHandler(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := Handler(path)
		if handler(c.Writer, c.Request) {
			c.JSON(200, Response{
				Alive: true,
			})
			c.Abort()
		}
	}
}
