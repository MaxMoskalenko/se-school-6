package gitsvc

import "context"

type Interface interface {
	FetchLatestReleaseTag(ctx context.Context, owner, repo string) (string, error)
	RepoExists(ctx context.Context, owner, repo string) (bool, error)
}
