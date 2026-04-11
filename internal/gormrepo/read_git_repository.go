package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) ReadGitRepository(ctx context.Context, params domain.ReadGitRepositoryParams) (*domain.GitRepository, error) {
	query := r.getDB(ctx)
	query, err := applyGitRepositoryFilters(query, params)
	if err != nil {
		return nil, err
	}

	var model gitRepositoryModel

	err = query.First(&model).Error
	if err == nil {
		return model.toDomain()
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if params.CreateIfNotExists == nil {
		return nil, gorm.ErrRecordNotFound
	}

	newModel := gitRepositoryModelFromDomain(params.CreateIfNotExists)
	if err := r.getDB(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(newModel).Error; err != nil {
		return nil, err
	}

	return newModel.toDomain()
}

func applyGitRepositoryFilters(query *gorm.DB, params domain.ReadGitRepositoryParams) (*gorm.DB, error) {
	if params.ByOwner == nil || params.ByName == nil {
		return nil, fmt.Errorf("missing filter: owner and name are required")
	}
	query = query.Where("owner = ?", strings.ToLower(*params.ByOwner))
	query = query.Where("name = ?", strings.ToLower(*params.ByName))
	return query, nil
}
