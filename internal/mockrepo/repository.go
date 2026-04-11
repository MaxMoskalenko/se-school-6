package mockrepo

import (
	"context"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/stretchr/testify/mock"
)

var _ domain.Repository = (*MockRepository)(nil)

type MockRepository struct {
	mock.Mock
}

func New() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return fn(ctx)
}

func (m *MockRepository) ReadUser(ctx context.Context, params domain.ReadUserParams) (*domain.User, error) {
	args := m.Called(ctx, params)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

func (m *MockRepository) ReadGitRepository(ctx context.Context, params domain.ReadGitRepositoryParams) (*domain.GitRepository, error) {
	args := m.Called(ctx, params)

	var repo *domain.GitRepository
	if args.Get(0) != nil {
		repo = args.Get(0).(*domain.GitRepository)
	}

	return repo, args.Error(1)
}

func (m *MockRepository) ReadGitRepositories(ctx context.Context, params domain.ReadGitRepositoriesParams) ([]*domain.GitRepository, error) {
	args := m.Called(ctx, params)

	var repos []*domain.GitRepository
	if args.Get(0) != nil {
		repos = args.Get(0).([]*domain.GitRepository)
	}

	return repos, args.Error(1)
}

func (m *MockRepository) SaveGitRepository(ctx context.Context, repo *domain.GitRepository) error {
	args := m.Called(ctx, repo)
	return args.Error(0)
}

func (m *MockRepository) SaveRepositorySubscription(ctx context.Context, subscription *domain.Subscription, params domain.SaveRepositorySubscriptionParams) error {
	args := m.Called(ctx, subscription, params)
	return args.Error(0)
}

func (m *MockRepository) ReadRepositorySubscription(ctx context.Context, params domain.ReadRepositorySubscriptionParams) (*domain.Subscription, error) {
	args := m.Called(ctx, params)

	var sub *domain.Subscription
	if args.Get(0) != nil {
		sub = args.Get(0).(*domain.Subscription)
	}

	return sub, args.Error(1)
}

func (m *MockRepository) ReadRepositorySubscriptions(ctx context.Context, params domain.ReadRepositorySubscriptionsParams) ([]*domain.Subscription, error) {
	args := m.Called(ctx, params)

	var subs []*domain.Subscription
	if args.Get(0) != nil {
		subs = args.Get(0).([]*domain.Subscription)
	}

	return subs, args.Error(1)
}
