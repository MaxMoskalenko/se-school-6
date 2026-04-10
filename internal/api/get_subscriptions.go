package api

import (
	"context"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
)

type GetSubscriptionsQuery struct {
	Email string
}

func (a *App) GetSubscriptions(ctx context.Context, q GetSubscriptionsQuery) ([]*domain.Subscription, error) {
	subs, err := a.repo.ReadRepositorySubscriptions(ctx, domain.ReadRepositorySubscriptionsParams{
		ByUserEmail:    &q.Email,
		WithRepository: true,
		WithDOITokens:  true,
	})
	if err != nil {
		return nil, err
	}

	return subs, nil
}
