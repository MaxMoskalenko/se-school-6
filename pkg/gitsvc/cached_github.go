package gitsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/MaxMoskalenko/se-school-6/pkg/cache"
)

type CachedGithubService struct {
	inner Interface
	cache cache.Interface
	ttl   time.Duration
}

func NewCachedGithubService(inner Interface, cache cache.Interface, ttl time.Duration) *CachedGithubService {
	return &CachedGithubService{inner: inner, cache: cache, ttl: ttl}
}

func (s *CachedGithubService) FetchLatestReleaseTag(ctx context.Context, owner, repo string) (string, error) {
	key := fmt.Sprintf("gitsvc:release:%s/%s", owner, repo)

	data, err := s.cache.Get(ctx, key)
	if err == nil {
		return string(data), nil
	}

	tag, err := s.inner.FetchLatestReleaseTag(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, key, []byte(tag), s.ttl)
	return tag, nil
}

func (s *CachedGithubService) RepoExists(ctx context.Context, owner, repo string) (bool, error) {
	key := fmt.Sprintf("gitsvc:exists:%s/%s", owner, repo)

	data, err := s.cache.Get(ctx, key)
	if err == nil {
		return string(data) == "1", nil
	}

	exists, err := s.inner.RepoExists(ctx, owner, repo)
	if err != nil {
		return false, err
	}

	val := "0"
	if exists {
		val = "1"
	}
	_ = s.cache.Set(ctx, key, []byte(val), s.ttl)
	return exists, nil
}
