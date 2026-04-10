package main

import (
	"database/sql"
	"embed"
	"log"

	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/MaxMoskalenko/se-school-6/internal/config"
	"github.com/MaxMoskalenko/se-school-6/internal/ginrouter"
	"github.com/MaxMoskalenko/se-school-6/internal/gormrepo"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

	app := api.NewApp(repo, api.Config{HostURL: cfg.Api.HostURL}, mailSvc)

	router, err := ginrouter.New(app, ginrouter.Config{
		Port: cfg.Router.Port,
	})
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func migrate(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(migrations)

	return goose.Up(db, "migrations/postgres")
}
