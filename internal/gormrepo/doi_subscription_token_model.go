package gormrepo

import (
	"fmt"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/google/uuid"
)

// using int16 for action to save space in the database
// also it's not domain model to avoid coupling domain and database
type doiSubscriptionTokenAction int16

const (
	doiSubscriptionTokenActionSubscribe   doiSubscriptionTokenAction = 0
	doiSubscriptionTokenActionUnsubscribe doiSubscriptionTokenAction = 1
)

func newDOISubscriptionTokenActionFromDomain(action domain.DOISubscriptionTokenAction) doiSubscriptionTokenAction {
	switch action {
	case domain.DOISubscriptionTokenActionSubscribe:
		return doiSubscriptionTokenActionSubscribe
	case domain.DOISubscriptionTokenActionUnsubscribe:
		return doiSubscriptionTokenActionUnsubscribe
	default:
		return doiSubscriptionTokenActionSubscribe // default to subscribe to avoid accidentally unsubscribing users, this should never happen anyway
	}
}

func (a doiSubscriptionTokenAction) toDomain() domain.DOISubscriptionTokenAction {
	switch a {
	case doiSubscriptionTokenActionSubscribe:
		return domain.DOISubscriptionTokenActionSubscribe
	case doiSubscriptionTokenActionUnsubscribe:
		return domain.DOISubscriptionTokenActionUnsubscribe
	default:
		return domain.DOISubscriptionTokenActionUnknown
	}
}

type doiSubscriptionTokenModel struct {
	ID             string                     `gorm:"column:id;primaryKey"`
	SubscriptionID string                     `gorm:"column:subscription_id;not null"`
	Action         doiSubscriptionTokenAction `gorm:"column:action;not null"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (doiSubscriptionTokenModel) TableName() string {
	return "doi_subscription_tokens"
}

func (m *doiSubscriptionTokenModel) toDomain() (*domain.DOISubscriptionToken, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("parse doi token id: %w", err)
	}

	action := doiSubscriptionTokenAction(m.Action).toDomain()
	return domain.NewDOISubscriptionToken(action).WithID(id), nil
}

func doiSubscriptionTokenModelFromDomain(t *domain.DOISubscriptionToken, subscriptionID string) *doiSubscriptionTokenModel {
	return &doiSubscriptionTokenModel{
		ID:             t.ID().String(),
		SubscriptionID: subscriptionID,
		Action:         newDOISubscriptionTokenActionFromDomain(t.Action()),
	}
}
