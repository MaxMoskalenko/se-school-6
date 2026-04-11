package scanner

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
)

func (a *App) Scan(ctx context.Context) error {
	repos, err := a.repo.ReadGitRepositories(ctx, domain.ReadGitRepositoriesParams{
		OnlyWithActiveSubscriptions: true,
		SortByLastCheckedAt:         true,
		WithSubscriptions:           true,
		WithUser:                    true,
	})
	if err != nil {
		return err
	}

	for _, repo := range repos {
		now := time.Now()
		latestTag, err := a.gitSvc.FetchLatestReleaseTag(ctx, repo.Owner(), repo.Name())
		if err != nil {
			if errors.Is(err, gitsvc.ErrRateLimited) {
				log.Printf("error: github rate limit reached, aborting scan")
				return err
			}
			log.Printf("error: failed to fetch latest tag for %s/%s: %v", repo.Owner(), repo.Name(), err)
			continue
		}

		isNewRelease := repo.LastSeenTag() == nil || *repo.LastSeenTag() != latestTag

		repo = repo.WithLastCheckedAt(&now)
		if isNewRelease {
			repo = repo.WithLastSeenTag(latestTag)
		}

		if err := a.repo.SaveGitRepository(ctx, repo); err != nil {
			log.Printf("error: failed to save git repository %s/%s: %v", repo.Owner(), repo.Name(), err)
			continue
		}

		if isNewRelease {
			log.Printf("new release detected for %s/%s: %s", repo.Owner(), repo.Name(), latestTag)
			a.notifySubscribers(ctx, repo, latestTag)
		}
	}

	return nil
}

func (a *App) notifySubscribers(ctx context.Context, repo *domain.GitRepository, tag string) {
	repoFullName := repo.Owner() + "/" + repo.Name()
	for _, sub := range repo.Subscriptions() {
		user := sub.User()
		if user == nil {
			continue
		}
		if err := a.mailSvc.SendNewReleaseEmail(ctx, mailsvc.NewReleaseEmailParams{
			Email:      user.Email(),
			Repo:       repoFullName,
			ReleaseTag: tag,
		}); err != nil {
			log.Printf("error: failed to send new release email to %s for %s: %v", user.Email(), repoFullName, err)
		}
	}
}
