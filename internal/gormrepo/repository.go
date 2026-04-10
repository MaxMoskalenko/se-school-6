package gormrepo

import (
	"fmt"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ domain.Repository = (*GormRepository)(nil)

type GormRepository struct {
	db *gorm.DB
}

func New(cfg GormConfig) (*GormRepository, error) {
	// postgres is hardcoded/default driver here, but it could be easily extended to support multiple drivers
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm connection: %w", err)
	}
	return &GormRepository{db: db}, nil
}

func NewFromDB(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}
