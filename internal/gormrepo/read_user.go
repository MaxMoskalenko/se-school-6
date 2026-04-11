package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) ReadUser(ctx context.Context, params domain.ReadUserParams) (*domain.User, error) {
	query := r.getDB(ctx)
	query, err := applyUserFilters(query, params)
	if err != nil {
		return nil, err
	}

	var model userModel

	err = query.First(&model).Error
	if err == nil {
		return model.toDomain()
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if params.CreateIfNotExists == nil {
		return nil, domain.ErrNotFound
	}

	newModel := userModelFromDomain(params.CreateIfNotExists)
	if err := r.getDB(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(newModel).Error; err != nil {
		return nil, err
	}

	return newModel.toDomain()
}

func applyUserFilters(query *gorm.DB, params domain.ReadUserParams) (*gorm.DB, error) {
	if params.ByEmail == nil {
		return nil, fmt.Errorf("missing filter: email is required")
	}
	query = query.Where("email = ?", *params.ByEmail)
	return query, nil
}
