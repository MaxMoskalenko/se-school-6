package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type HandlerFactory struct {
	app *api.App
	cfg Config
}

func NewHandlerFactory(app *api.App, cfg Config) *HandlerFactory {
	return &HandlerFactory{app: app, cfg: cfg}
}

func (hf *HandlerFactory) Handler(handler func(app *api.App, cfg Config, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(hf.app, hf.cfg, c)
	}
}
