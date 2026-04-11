package api

import (
	"context"
	"log"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
)

type GetSubscriptionsQuery struct {
	Email string
}

func (a *App) GetSubscriptions(ctx context.Context, q GetSubscriptionsQuery) ([]*domain.Subscription, error) {
	subs, err := a.repo.ReadRepositorySubscriptions(ctx, domain.ReadRepositorySubscriptionsParams{
		ByUserEmail:    &q.Email,
		WithRepository: true,
		WithUser:       true,
	})
	if err != nil {
		log.Printf("error: failed to read subscriptions for email=%s: %v", q.Email, err)
		return nil, domain.ErrInternal
	}

	return subs, nil
}
