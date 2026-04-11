package domain

import "context"

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	Transactor

	ReadUser(ctx context.Context, params ReadUserParams) (*User, error)

	ReadGitRepository(ctx context.Context, params ReadGitRepositoryParams) (*GitRepository, error)
	ReadGitRepositories(ctx context.Context, params ReadGitRepositoriesParams) ([]*GitRepository, error)
	SaveGitRepository(ctx context.Context, repo *GitRepository) error

	SaveRepositorySubscription(ctx context.Context, subscription *Subscription, params SaveRepositorySubscriptionParams) error
	ReadRepositorySubscription(ctx context.Context, params ReadRepositorySubscriptionParams) (*Subscription, error)
	ReadRepositorySubscriptions(ctx context.Context, params ReadRepositorySubscriptionsParams) ([]*Subscription, error)
}

type ReadUserParams struct {
	ByEmail *string

	CreateIfNotExists *User
}

type ReadGitRepositoriesParams struct {
	OnlyWithActiveSubscriptions bool
	SortByLastCheckedAt         bool

	WithSubscriptions bool
	WithUser          bool
}

type ReadGitRepositoryParams struct {
	ByOwner *string
	ByName  *string

	CreateIfNotExists *GitRepository
}

type SaveRepositorySubscriptionParams struct {
	SaveDOITokens bool
}

type ReadRepositorySubscriptionParams struct {
	ByDOIToken          *string
	ByUserID            *string
	ByRepositoryID      *string
	OnlyNonUnsubscribed bool

	WithDOITokens  bool
	WithUser       bool
	WithRepository bool
}

type ReadRepositorySubscriptionsParams struct {
	ByUserEmail *string
	OnlyActive  bool

	WithDOITokens  bool
	WithUser       bool
	WithRepository bool
}
