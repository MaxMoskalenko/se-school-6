package gormrepo

import (
	"fmt"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/google/uuid"
)

type repositorySubscriptionModel struct {
	ID             string     `gorm:"column:id;primaryKey"`
	UserID         string     `gorm:"column:user_id;not null"`
	RepositoryID   string     `gorm:"column:repository_id;not null"`
	ConfirmedAt    *time.Time `gorm:"column:confirmed_at"`
	UnsubscribedAt *time.Time `gorm:"column:unsubscribed_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	User               *userModel                  `gorm:"foreignKey:UserID;references:ID"`
	Repository         *gitRepositoryModel         `gorm:"foreignKey:RepositoryID;references:ID"`
	SubscriptionTokens []doiSubscriptionTokenModel `gorm:"foreignKey:SubscriptionID;references:ID"`
}

func (repositorySubscriptionModel) TableName() string {
	return "repository_subscriptions"
}

func (m *repositorySubscriptionModel) toDomain() (*domain.Subscription, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("parse subscription id: %w", err)
	}

	sub := domain.NewSubscription().
		WithID(id).
		WithConfirmedAt(m.ConfirmedAt).
		WithUnsubscribedAt(m.UnsubscribedAt)

	if m.User != nil {
		u, err := m.User.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map user: %w", err)
		}
		sub = sub.WithUser(u)
	}

	if m.Repository != nil {
		r, err := m.Repository.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map repository: %w", err)
		}
		sub = sub.WithGitRepository(r)
	}

	for i := range m.SubscriptionTokens {
		t, err := m.SubscriptionTokens[i].toDomain()
		if err != nil {
			return nil, fmt.Errorf("map doi token: %w", err)
		}
		sub = sub.WithDOISubscriptionToken(t)
	}

	return sub, nil
}

func subscriptionModelFromDomain(s *domain.Subscription) *repositorySubscriptionModel {
	model := &repositorySubscriptionModel{
		ID:             s.ID().String(),
		ConfirmedAt:    s.ConfirmedAt(),
		UnsubscribedAt: s.UnsubscribedAt(),
	}

	if u := s.User(); u != nil {
		model.UserID = u.ID().String()
	}

	if r := s.GitRepository(); r != nil {
		model.RepositoryID = r.ID().String()
	}

	return model
}
