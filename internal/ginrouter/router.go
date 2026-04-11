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
	factory := NewHandlerFactory(app, cfg)

	if err := bindvalidator.Register(); err != nil {
		return Router{}, err
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	apiGroup := r.Group("/api")
	apiGroup.POST("/auth", factory.Handler(postAuthHandler))

	authorized := apiGroup.Group("")

	if cfg.ValidateAuthEmail {
		authorized.Use(authMiddleware(cfg.JWTSecret))
	}

	authorized.POST("/subscribe", factory.Handler(postSubscribeHandler))
	authorized.GET("/confirm/:token", factory.Handler(getConfirmHandler))
	authorized.GET("/unsubscribe/:token", factory.Handler(getUnsubscribeHandler))
	authorized.GET("/subscriptions", factory.Handler(getSubscriptionsHandler))

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
