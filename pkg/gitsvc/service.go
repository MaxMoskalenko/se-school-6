package githubsvc

import (
	"context"

	"github.com/google/go-github/v72/github"
)

type GithubService struct {
	client *github.Client
	cfg    GithubConfig
}

func NewGithubService(cfg GithubConfig) *GithubService {
	client := github.NewClient(nil)

	if cfg.AuthToken != nil {
		client = client.WithAuthToken(*cfg.AuthToken)
	}

	return &GithubService{
		client: client,
		cfg:    cfg,
	}
}

func (s *GithubService) FetchLatestReleaseTag(ctx context.Context, owner, repo string) (string, error) {
	latest, _, err := s.client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", err
	}
	// Do something with the latest release
	_ = latest
	return latest.GetTagName(), nil
}
