package gormrepo

import (
	"context"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
)

func (r *GormRepository) SaveRepositorySubscription(ctx context.Context, subscription *domain.Subscription, params domain.SaveRepositorySubscriptionParams) error {
	model := subscriptionModelFromDomain(subscription)

	return r.WithTransaction(ctx, func(ctx context.Context) error {
		columnsToUpdate := defineSubscriptionUpdateColumns(subscription)

		if len(columnsToUpdate) > 0 {
			if err := r.getDB(ctx).
				Model(&repositorySubscriptionModel{ID: model.ID}).
				Select(columnsToUpdate).
				Updates(model).Error; err != nil {
				return err
			}
			return nil
		}

		if err := r.getDB(ctx).Create(model).Error; err != nil {
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

func defineSubscriptionUpdateColumns(subscription *domain.Subscription) []string {
	var columns []string
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
