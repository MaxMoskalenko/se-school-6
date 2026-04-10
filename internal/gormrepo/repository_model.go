package gormrepo

import (
	"fmt"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/google/uuid"
)

type gitRepositoryModel struct {
	ID          string  `gorm:"column:id;primaryKey"`
	Name        string  `gorm:"column:name;not null"`
	Owner       string  `gorm:"column:owner;not null"`
	LastSeenTag *string `gorm:"column:last_seen_tag"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (gitRepositoryModel) TableName() string {
	return "git_repositories"
}

func (m *gitRepositoryModel) toDomain() (*domain.GitRepository, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("parse repository id: %w", err)
	}

	r := domain.NewGitRepository(m.Owner, m.Name).WithID(id)
	if m.LastSeenTag != nil {
		r = r.WithLastSeenTag(*m.LastSeenTag)
	}
	return r, nil
}

func gitRepositoryModelFromDomain(r *domain.GitRepository) *gitRepositoryModel {
	model := &gitRepositoryModel{
		ID:    r.ID().String(),
		Name:  r.Name(),
		Owner: r.Owner(),
	}
	if tag := r.LastSeenTag(); tag != nil {
		model.LastSeenTag = tag
	}
	return model
}
