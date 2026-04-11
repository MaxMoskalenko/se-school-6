package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os/signal"
	"syscall"

	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/MaxMoskalenko/se-school-6/internal/config"
	"github.com/MaxMoskalenko/se-school-6/internal/ginrouter"
	"github.com/MaxMoskalenko/se-school-6/internal/gormrepo"
	"github.com/MaxMoskalenko/se-school-6/internal/scanner"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
)

//go:embed migrations/postgres/*.sql
var migrations embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbCfg := gormrepo.GormConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}

	if err := migrate(dbCfg.DSN()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo, err := gormrepo.New(dbCfg)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	mailSvc := mailsvc.NewPostmark(mailsvc.PostmarkConfig{
		ServerToken:                cfg.Postmark.ServerToken,
		AccountToken:               cfg.Postmark.AccountToken,
		SenderEmail:                cfg.Postmark.SenderEmail,
		SubscribeRequestTemplateID: cfg.Postmark.SubscribeRequestTemplateID,
		NewReleaseTemplateID:       cfg.Postmark.NewReleaseTemplateID,
	})

	gitSvc := gitsvc.NewGithubService(gitsvc.GithubConfig{
		AuthToken: cfg.Github.AuthToken,
	})

	app := api.NewApp(repo, api.Config{HostURL: cfg.Api.HostURL, JWTSecret: cfg.Api.JWTSecret}, mailSvc, gitSvc)

	router, err := ginrouter.New(app, ginrouter.Config{
		Port:              cfg.Router.Port,
		JWTSecret:         cfg.Api.JWTSecret,
		ValidateAuthEmail: cfg.Api.ValidateAuthEmail,
	})
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	sc := scanner.NewApp(repo, scanner.Config{
		Interval: cfg.Scanner.Interval,
	}, gitSvc, mailSvc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return router.Run(ctx)
	})

	g.Go(func() error {
		return sc.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("application stopped: %v", err)
	}
}

func migrate(dsn string) (err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	goose.SetBaseFS(migrations)

	return goose.Up(db, "migrations/postgres")
}
