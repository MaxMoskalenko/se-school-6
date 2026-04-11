package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/internal/mockrepo"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSubscriptions_Success(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	subs := []*domain.Subscription{
		domain.NewSubscription().WithNewID(),
		domain.NewSubscription().WithNewID(),
	}

	repo.On("ReadRepositorySubscriptions", mock.Anything, mock.Anything).Return(subs, nil)

	result, err := app.GetSubscriptions(context.Background(), GetSubscriptionsQuery{
		Email: "test@example.com",
	})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetSubscriptions_Empty(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscriptions", mock.Anything, mock.Anything).Return([]*domain.Subscription{}, nil)

	result, err := app.GetSubscriptions(context.Background(), GetSubscriptionsQuery{
		Email: "test@example.com",
	})

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetSubscriptions_DBError(t *testing.T) {
	repo := mockrepo.New()
	app := newTestApp(repo, mailsvc.NewMock(), gitsvc.NewMock())

	repo.On("ReadRepositorySubscriptions", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	result, err := app.GetSubscriptions(context.Background(), GetSubscriptionsQuery{
		Email: "test@example.com",
	})

	assert.ErrorIs(t, err, domain.ErrInternal)
	assert.Nil(t, result)
}
