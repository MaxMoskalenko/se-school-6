package api

import (
	"context"
	"errors"
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
		ByDOIToken:     &cmd.Token,
		WithDOITokens:  true,
		WithUser:       true,
		WithRepository: true,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.NewError(http.StatusNotFound, domain.ErrNotFound)
		}
		log.Printf("error: failed to read subscription by DOI token: %v", err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	if sub.DOITokenAction(cmd.Token) != domain.DOISubscriptionTokenActionUnsubscribe {
		return domain.NewError(http.StatusBadRequest, domain.ErrInvalidToken)
	}

	if !sub.IsActive() {
		return domain.NewError(http.StatusConflict, domain.ErrNotActive)
	}

	now := time.Now()
	sub = sub.WithUnsubscribedAt(&now)
	if err := a.repo.SaveRepositorySubscription(ctx, sub, domain.SaveRepositorySubscriptionParams{}); err != nil {
		log.Printf("error: failed to save unsubscribed subscription: %v", err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	return nil
}
