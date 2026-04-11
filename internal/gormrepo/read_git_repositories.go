package gormrepo

import (
	"context"
	"fmt"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
)

func (r *GormRepository) ReadGitRepositories(ctx context.Context, params domain.ReadGitRepositoriesParams) ([]*domain.GitRepository, error) {
	query := r.getDB(ctx).Model(&gitRepositoryModel{})
	query = applyGitRepositoriesJoins(query, params)
	query = applyGitRepositoriesFilters(query, params)
	query = applyGitRepositoriesOrder(query, params)

	var models []gitRepositoryModel
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	repos := make([]*domain.GitRepository, 0, len(models))
	for i := range models {
		repo, err := models[i].toDomain()
		if err != nil {
			return nil, fmt.Errorf("map git repository: %w", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func applyGitRepositoriesJoins(query *gorm.DB, params domain.ReadGitRepositoriesParams) *gorm.DB {
	if params.WithSubscriptions {
		preload := query.Preload("Subscriptions", "confirmed_at IS NOT NULL AND unsubscribed_at IS NULL")
		if params.WithUser {
			preload = preload.Preload("Subscriptions.User")
		}
		return preload
	}
	return query
}

func applyGitRepositoriesFilters(query *gorm.DB, params domain.ReadGitRepositoriesParams) *gorm.DB {
	if params.OnlyWithActiveSubscriptions {
		query = query.Where("EXISTS (SELECT 1 FROM repository_subscriptions rs WHERE rs.repository_id = git_repositories.id AND rs.confirmed_at IS NOT NULL AND rs.unsubscribed_at IS NULL)")
	}
	return query
}

func applyGitRepositoriesOrder(query *gorm.DB, params domain.ReadGitRepositoriesParams) *gorm.DB {
	if params.SortByLastCheckedAt {
		query = query.Order("last_checked_at ASC NULLS FIRST")
	}
	return query
}
