package gormrepo

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func (r *GormRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func (r *GormRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}
