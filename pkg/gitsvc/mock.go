package gitsvc

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

func (m *Mock) FetchLatestReleaseTag(ctx context.Context, owner, repo string) (string, error) {
	args := m.Called(ctx, owner, repo)
	return args.String(0), args.Error(1)
}

func (m *Mock) RepoExists(ctx context.Context, owner, repo string) (bool, error) {
	args := m.Called(ctx, owner, repo)
	return args.Bool(0), args.Error(1)
}
