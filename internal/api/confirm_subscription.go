package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
)

type ConfirmSubscriptionCommand struct {
	Token string
}

func (a *App) ConfirmSubscription(ctx context.Context, cmd ConfirmSubscriptionCommand) *domain.Error {
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

	if sub.DOITokenAction(cmd.Token) != domain.DOISubscriptionTokenActionSubscribe {
		return domain.NewError(http.StatusBadRequest, domain.ErrInvalidToken)
	}

	if sub.IsConfirmed() {
		return domain.NewError(http.StatusConflict, domain.ErrAlreadyConfirmed)
	}

	now := time.Now()
	sub = sub.WithConfirmedAt(&now)
	if err := a.repo.SaveRepositorySubscription(ctx, sub, domain.SaveRepositorySubscriptionParams{}); err != nil {
		log.Printf("error: failed to save confirmed subscription: %v", err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	return nil
}
