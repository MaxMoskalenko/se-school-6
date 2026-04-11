package gormrepo

import (
	"context"
	"errors"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) SaveRepositorySubscription(ctx context.Context, subscription *domain.Subscription, params domain.SaveRepositorySubscriptionParams) error {
	model := subscriptionModelFromDomain(subscription)

	// if confirmed_at or unsubscribed_at is set, updates the existing record,
	// otherwise returns an error if the record already exists (to prevent duplicate subscriptions)
	columnsToUpdate := r.defineSubscriptionUpsertColumns(subscription)

	return r.WithTransaction(ctx, func(ctx context.Context) error {
		query := r.getDB(ctx)

		if len(columnsToUpdate) > 0 {
			query = query.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "repository_id"}},
				DoUpdates: clause.AssignmentColumns(columnsToUpdate),
			})
		}

		if err := query.Create(model).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrAlreadySubscribed
			}
			return err
		}

		if params.SaveDOITokens {
			for _, token := range subscription.DOISubscriptionTokens() {
				tokenModel := doiSubscriptionTokenModelFromDomain(token, model.ID)
				if err := r.getDB(ctx).Create(tokenModel).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *GormRepository) defineSubscriptionUpsertColumns(subscription *domain.Subscription) []string {
	columns := []string{}
	if subscription.ConfirmedAt() != nil {
		columns = append(columns, "confirmed_at")
	}
	if subscription.UnsubscribedAt() != nil {
		columns = append(columns, "unsubscribed_at")
	}

	if len(columns) > 0 {
		columns = append(columns, "updated_at")
	}

	return columns
}
