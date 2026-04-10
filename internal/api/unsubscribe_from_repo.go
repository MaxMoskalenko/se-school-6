package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
)

type UnsubscribeCommand struct {
	Token string
}

func (a *App) UnsubscribeFromRepo(ctx context.Context, cmd UnsubscribeCommand) *domain.Error {
	sub, err := a.repo.ReadRepositorySubscription(ctx, domain.ReadRepositorySubscriptionParams{
		ByDOIToken:    &cmd.Token,
		WithDOITokens:  true,
		WithUser:       true,
		WithRepository: true,
	})
	if err != nil {
		log.Printf("error: failed to read subscription by DOI token: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	if sub.DOITokenAction(cmd.Token) != domain.DOISubscriptionTokenActionUnsubscribe {
		log.Printf("error: invalid DOI token action for unsubscribe, token=%s", cmd.Token)
		return domain.NewError(http.StatusBadRequest, errInvalidToken)
	}

	if !sub.IsActive() {
		log.Printf("error: subscription is not active, token=%s", cmd.Token)
		return domain.NewError(http.StatusConflict, errNotActive)
	}

	now := time.Now()
	sub = sub.WithUnsubscribedAt(&now)
	if err := a.repo.SaveRepositorySubscription(ctx, sub, domain.SaveRepositorySubscriptionParams{}); err != nil {
		log.Printf("error: failed to save unsubscribed subscription: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	return nil
}
