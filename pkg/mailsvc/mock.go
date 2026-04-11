package mailsvc

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Interface = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) SendSubscribeRequestEmail(ctx context.Context, params SubscribeRequestParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *Mock) SendNewReleaseEmail(ctx context.Context, params NewReleaseEmailParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}
