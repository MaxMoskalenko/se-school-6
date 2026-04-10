package gormrepo

import (
	"fmt"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/google/uuid"
)

type userModel struct {
	ID    string `gorm:"column:id;primaryKey"`
	Email string `gorm:"column:email;not null;uniqueIndex"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (userModel) TableName() string {
	return "users"
}

func (m *userModel) toDomain() (*domain.User, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	u := domain.NewUser(m.Email).WithID(id)
	return u, nil
}

func userModelFromDomain(u *domain.User) *userModel {
	return &userModel{
		ID:    u.ID().String(),
		Email: u.Email(),
	}
}
