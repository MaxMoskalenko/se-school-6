package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/MaxMoskalenko/se-school-6/pkg/bindvalidator"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	cfg Config
}

func New(app *api.App, cfg Config) (Router, error) {
	factory := NewHandlerFactory(app)

	// Register custom validation rules
	if err := bindvalidator.Register(); err != nil {
		return Router{}, err
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	// Public routes
	r.GET("/auth", getAuthHandler)

	// API routes
	apiGroup := r.Group("/api")
	apiGroup.POST("/subscribe", factory.Handler(postSubscribeHandler))
	apiGroup.GET("/confirm/:token", factory.Handler(getConfirmHandler))
	apiGroup.GET("/unsubscribe/:token", factory.Handler(getUnsubscribeHandler))
	apiGroup.GET("/subscriptions", factory.Handler(getSubscriptionsHandler))

	return Router{
		Engine: r,
		cfg:    cfg,
	}, nil
}

func (r Router) Run() error {
	return r.Engine.Run(":" + r.cfg.Port)
}
