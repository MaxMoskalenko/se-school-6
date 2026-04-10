package api

import (
	"context"
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
		ByDOIToken:    &cmd.Token,
		WithDOITokens:  true,
		WithUser:       true,
		WithRepository: true,
	})
	if err != nil {
		log.Printf("error: failed to read subscription by DOI token: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	if sub.DOITokenAction(cmd.Token) != domain.DOISubscriptionTokenActionSubscribe {
		log.Printf("error: invalid DOI token action for confirmation, token=%s", cmd.Token)
		return domain.NewError(http.StatusBadRequest, errInvalidToken)
	}

	if sub.IsConfirmed() {
		log.Printf("error: subscription already confirmed, token=%s", cmd.Token)
		return domain.NewError(http.StatusConflict, errAlreadyConfirmed)
	}

	now := time.Now()
	sub = sub.WithConfirmedAt(&now)
	if err := a.repo.SaveRepositorySubscription(ctx, sub, domain.SaveRepositorySubscriptionParams{}); err != nil {
		log.Printf("error: failed to save confirmed subscription: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	return nil
}
