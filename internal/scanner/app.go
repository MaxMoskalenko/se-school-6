package scanner

import (
	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
)

type App struct {
	repo    domain.Repository
	cfg     Config
	gitSvc  gitsvc.Interface
	mailSvc mailsvc.Interface
}

func NewApp(repo domain.Repository, cfg Config, gitSvc gitsvc.Interface, mailSvc mailsvc.Interface) *App {
	return &App{
		repo:    repo,
		cfg:     cfg,
		gitSvc:  gitSvc,
		mailSvc: mailSvc,
	}
}
