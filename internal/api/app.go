package api

import (
	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
)

type App struct {
	repo domain.Repository
	cfg  Config

	mailSvc mailsvc.Interface
}

func NewApp(repo domain.Repository, cfg Config, mailSvc mailsvc.Interface) *App {
	return &App{
		repo:    repo,
		cfg:     cfg,
		mailSvc: mailSvc,
	}
}
