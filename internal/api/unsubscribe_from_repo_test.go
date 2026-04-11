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

func TestUnsubscribeFromRepo_Success(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	unsubTokenID := uuid.New()
	confirmedAt := time.Now()
	sub := buildSubscription(uuid.New(), unsubTokenID, &confirmedAt)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: unsubTokenID.String(),
	})

	assert.Nil(t, dErr)
	repo.AssertExpectations(t)
}

func TestUnsubscribeFromRepo_TokenNotFound(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: uuid.New().String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusNotFound, dErr.Code())
}

func TestUnsubscribeFromRepo_DBError(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: uuid.New().String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}

func TestUnsubscribeFromRepo_InvalidToken_WrongAction(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	subscribeTokenID := uuid.New()
	confirmedAt := time.Now()
	sub := buildSubscription(subscribeTokenID, uuid.New(), &confirmedAt)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)

	// use the subscribe token for an unsubscribe request
	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: subscribeTokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusBadRequest, dErr.Code())
}

func TestUnsubscribeFromRepo_NotActive_Unconfirmed(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	unsubTokenID := uuid.New()
	// not confirmed yet — IsActive() returns false
	sub := buildSubscription(uuid.New(), unsubTokenID, nil)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: unsubTokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusConflict, dErr.Code())
}

func TestUnsubscribeFromRepo_NotActive_AlreadyUnsubscribed(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	unsubTokenID := uuid.New()
	confirmedAt := time.Now()
	unsubscribedAt := time.Now()
	sub := buildSubscription(uuid.New(), unsubTokenID, &confirmedAt).
		WithUnsubscribedAt(&unsubscribedAt)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: unsubTokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusConflict, dErr.Code())
}

func TestUnsubscribeFromRepo_SaveFails(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	unsubTokenID := uuid.New()
	confirmedAt := time.Now()
	sub := buildSubscription(uuid.New(), unsubTokenID, &confirmedAt)

	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(sub, nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("save error"))

	dErr := app.UnsubscribeFromRepo(context.Background(), UnsubscribeCommand{
		Token: unsubTokenID.String(),
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}
