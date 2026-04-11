package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Api      Api      `mapstructure:"api"`
	Router   Router   `mapstructure:"router"`
	Database Database `mapstructure:"database"`
	Postmark Postmark `mapstructure:"postmark"`
	Github   Github   `mapstructure:"github"`
	Scanner  Scanner  `mapstructure:"scanner"`
}

type Router struct {
	Port string `mapstructure:"port"`
}

type Api struct {
	HostURL string `mapstructure:"host_url"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type Postmark struct {
	ServerToken  string `mapstructure:"server_token"`
	AccountToken string `mapstructure:"account_token"`
	SenderEmail  string `mapstructure:"sender_email"`

	SubscribeRequestTemplateID int64 `mapstructure:"subscribe_request_template_id"`
	NewReleaseTemplateID       int64 `mapstructure:"new_release_template_id"`
}

type Github struct {
	AuthToken *string `mapstructure:"auth_token"`
}

type Scanner struct {
	Interval time.Duration
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	cfg := &Config{
		Api: Api{
			HostURL: viper.GetString("API_HOST_URL"),
		},
		Router: Router{
			Port: viper.GetString("ROUTER_PORT"),
		},
		Database: Database{
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetInt("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Name:     viper.GetString("DATABASE_NAME"),
			SSLMode:  viper.GetString("DATABASE_SSL_MODE"),
		},
		Postmark: Postmark{
			ServerToken:                viper.GetString("POSTMARK_SERVER_TOKEN"),
			AccountToken:               viper.GetString("POSTMARK_ACCOUNT_TOKEN"),
			SenderEmail:                viper.GetString("POSTMARK_SENDER_EMAIL"),
			SubscribeRequestTemplateID: viper.GetInt64("POSTMARK_SUBSCRIBE_REQUEST_TEMPLATE_ID"),
			NewReleaseTemplateID:       viper.GetInt64("POSTMARK_NEW_RELEASE_TEMPLATE_ID"),
		},
	}

	scanInterval, err := time.ParseDuration(viper.GetString("SCANNER_INTERVAL"))
	if err != nil {
		return nil, fmt.Errorf("parse scanner interval: %w", err)
	}
	cfg.Scanner = Scanner{
		Interval: scanInterval,
	}

	// this is one is optional, so parsing is a little bit different
	if token := viper.GetString("GITHUB_AUTH_TOKEN"); token != "" {
		cfg.Github.AuthToken = &token
	}

	return cfg, nil
}
