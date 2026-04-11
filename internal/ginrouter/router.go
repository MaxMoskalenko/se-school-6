package ginrouter

import (
	"context"
	"net/http"

	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/MaxMoskalenko/se-school-6/pkg/bindvalidator"
	"github.com/gin-gonic/gin"
)

type Router struct {
	server *http.Server
}

func New(app *api.App, cfg Config) (Router, error) {
	factory := NewHandlerFactory(app)

	if err := bindvalidator.Register(); err != nil {
		return Router{}, err
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	apiGroup := r.Group("/api")
	apiGroup.POST("/subscribe", factory.Handler(postSubscribeHandler))
	apiGroup.GET("/confirm/:token", factory.Handler(getConfirmHandler))
	apiGroup.GET("/unsubscribe/:token", factory.Handler(getUnsubscribeHandler))
	apiGroup.GET("/subscriptions", factory.Handler(getSubscriptionsHandler))

	return Router{
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: r,
		},
	}, nil
}

func (r Router) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return r.server.Shutdown(context.Background())
	}
}
