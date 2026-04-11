package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/internal/mockrepo"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildSubscription(subscribeTokenID, unsubscribeTokenID uuid.UUID, confirmedAt *time.Time) *domain.Subscription {
	sub := domain.NewSubscription().WithNewID()

	subscribeTok := domain.NewDOISubscriptionToken(domain.DOISubscriptionTokenActionSubscribe).WithID(subscribeTokenID)
	unsubscribeTok := domain.NewDOISubscriptionToken(domain.DOISubscriptionTokenActionUnsubscribe).WithID(unsubscribeTokenID)

	sub = sub.
		WithDOISubscriptionToken(subscribeTok).
		WithDOISubscriptionToken(unsubscribeTok).
		WithConfirmedAt(confirmedAt)

	return sub
}

func TestConfirmSubscription_Success(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	tokenID := uuid.New()
	sub := buildSubscription(tokenID, uuid.New(), nil)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: tokenID.String(),
	})

	assert.Nil(t, dErr)
	repo.AssertExpectations(t)
}

func TestConfirmSubscription_TokenNotFound(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)

	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: uuid.New().String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusNotFound, dErr.Code())
}

func TestConfirmSubscription_DBError(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: uuid.New().String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}

func TestConfirmSubscription_InvalidToken_WrongAction(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	unsubscribeTokenID := uuid.New()
	sub := buildSubscription(uuid.New(), unsubscribeTokenID, nil)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)

	// use the unsubscribe token for a confirm request
	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: unsubscribeTokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusBadRequest, dErr.Code())
}

func TestConfirmSubscription_AlreadyConfirmed(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	tokenID := uuid.New()
	confirmedAt := time.Now()
	sub := buildSubscription(tokenID, uuid.New(), &confirmedAt)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)

	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: tokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusConflict, dErr.Code())
}

func TestConfirmSubscription_SaveFails(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	tokenID := uuid.New()
	sub := buildSubscription(tokenID, uuid.New(), nil)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("save error"))

	dErr := app.ConfirmSubscription(context.Background(), ConfirmSubscriptionCommand{
		Token: tokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}
