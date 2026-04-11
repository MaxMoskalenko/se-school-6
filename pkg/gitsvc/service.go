package gitsvc

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v72/github"
)

var ErrRateLimited = errors.New("github API rate limit reached")

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
		return "", wrapRateLimitErr(err)
	}
	return latest.GetTagName(), nil
}

func (s *GithubService) RepoExists(ctx context.Context, owner, repo string) (bool, error) {
	_, resp, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, wrapRateLimitErr(err)
	}
	return true, nil
}

func wrapRateLimitErr(err error) error {
	var rateLimitErr *github.RateLimitError
	var abuseErr *github.AbuseRateLimitError
	if errors.As(err, &rateLimitErr) || errors.As(err, &abuseErr) {
		return ErrRateLimited
	}
	return err
}
