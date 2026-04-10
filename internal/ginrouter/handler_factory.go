package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type HandlerFactory struct {
	app *api.App
}

func NewHandlerFactory(app *api.App) *HandlerFactory {
	return &HandlerFactory{app: app}
}

func (hf *HandlerFactory) Handler(handler func(app *api.App, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(hf.app, c)
	}
}
