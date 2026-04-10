package gormrepo

import (
	"context"
	"fmt"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
)

func (r *GormRepository) ReadRepositorySubscriptions(ctx context.Context, params domain.ReadRepositorySubscriptionsParams) ([]*domain.Subscription, error) {
	query := r.getDB(ctx)
	query = applyRepositorySubscriptionsJoins(query, params)
	query, err := applyRepositorySubscriptionsFilters(query, params)
	if err != nil {
		return nil, err
	}

	var models []repositorySubscriptionModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	subs := make([]*domain.Subscription, 0, len(models))
	for i := range models {
		sub, err := models[i].toDomain()
		if err != nil {
			return nil, fmt.Errorf("map subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func applyRepositorySubscriptionsJoins(query *gorm.DB, params domain.ReadRepositorySubscriptionsParams) *gorm.DB {
	if params.WithUser || params.ByUserEmail != nil {
		query = query.Joins("User")
	}
	if params.WithRepository {
		query = query.Joins("Repository")
	}
	if params.WithDOITokens {
		query = query.Preload("SubscriptionTokens")
	}
	return query
}

func applyRepositorySubscriptionsFilters(query *gorm.DB, params domain.ReadRepositorySubscriptionsParams) (*gorm.DB, error) {
	if params.ByUserEmail == nil {
		return nil, fmt.Errorf("missing filter: user email is required")
	}
	query = query.Where("\"User\".email = ?", *params.ByUserEmail)
	return query, nil
}
