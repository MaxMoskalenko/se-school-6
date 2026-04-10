package gormrepo

import (
	"context"
	"fmt"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
)

func (r *GormRepository) ReadRepositorySubscription(ctx context.Context, params domain.ReadRepositorySubscriptionParams) (*domain.Subscription, error) {
	query := r.getDB(ctx)
	query = applyRepositorySubscriptionJoins(query, params)
	query = applyRepositorySubscriptionFilters(query, params)

	var model repositorySubscriptionModel
	if err := query.First(&model).Error; err != nil {
		return nil, err
	}

	sub, err := model.toDomain()
	if err != nil {
		return nil, fmt.Errorf("map subscription: %w", err)
	}

	return sub, nil
}

func applyRepositorySubscriptionJoins(query *gorm.DB, params domain.ReadRepositorySubscriptionParams) *gorm.DB {
	if params.WithUser {
		query = query.Joins("User")
	}
	if params.WithRepository {
		query = query.Joins("Repository")
	}
	if params.WithDOITokens {
		query = query.Preload("SubscriptionTokens")
	}
	if params.ByDOIToken != nil {
		query = query.Joins("JOIN doi_subscription_tokens ON doi_subscription_tokens.subscription_id = repository_subscriptions.id")
	}
	return query
}

func applyRepositorySubscriptionFilters(query *gorm.DB, params domain.ReadRepositorySubscriptionParams) *gorm.DB {
	if params.ByDOIToken != nil {
		query = query.Where("doi_subscription_tokens.id = ?", *params.ByDOIToken)
	}
	return query
}
