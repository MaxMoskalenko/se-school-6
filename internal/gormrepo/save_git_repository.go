package gormrepo

import (
	"context"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) SaveGitRepository(ctx context.Context, repo *domain.GitRepository) error {
	model := gitRepositoryModelFromDomain(repo)
	return r.getDB(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "owner"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_seen_tag", "last_checked_at", "updated_at"}),
		}).
		Create(model).Error
}
