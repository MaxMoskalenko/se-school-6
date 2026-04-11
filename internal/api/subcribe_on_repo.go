package api

import (
	"context"
	"log"
	"net/http"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
)

type SubscribeOnRepoCommand struct {
	Email     string
	RepoOwner string
	RepoName  string
}

func (a *App) SubscribeOnRepo(ctx context.Context, cmd SubscribeOnRepoCommand) *domain.Error {
	userToCreate := domain.NewUser(cmd.Email).WithNewID()

	// there are three possible db writes here, and they are not in a transaction,
	// but that's ok, since atomicity is not required, and in the worst case we will
	// have some orphaned users or repositories without subscription

	user, err := a.repo.ReadUser(ctx, domain.ReadUserParams{
		ByEmail:           &cmd.Email,
		CreateIfNotExists: userToCreate,
	})
	if err != nil {
		log.Printf("error: failed to read or create user email=%s: %v", cmd.Email, err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	gitRepo, err := a.repo.ReadGitRepository(ctx, domain.ReadGitRepositoryParams{
		ByOwner: &cmd.RepoOwner,
		ByName:  &cmd.RepoName,
	})
	if err != nil {
		exists, ghErr := a.gitSvc.RepoExists(ctx, cmd.RepoOwner, cmd.RepoName)
		if ghErr != nil {
			log.Printf("error: failed to check repo on github owner=%s name=%s: %v", cmd.RepoOwner, cmd.RepoName, ghErr)
			return domain.NewError(http.StatusInternalServerError, ghErr)
		}
		if !exists {
			log.Printf("error: repo not found on github owner=%s name=%s", cmd.RepoOwner, cmd.RepoName)
			return domain.NewError(http.StatusNotFound, errRepoNotFound)
		}

		gitRepoToCreate := domain.NewGitRepository(cmd.RepoOwner, cmd.RepoName).WithNewID()
		gitRepo, err = a.repo.ReadGitRepository(ctx, domain.ReadGitRepositoryParams{
			ByOwner:           &cmd.RepoOwner,
			ByName:            &cmd.RepoName,
			CreateIfNotExists: gitRepoToCreate,
		})
		if err != nil {
			log.Printf("error: failed to create git repository owner=%s name=%s: %v", cmd.RepoOwner, cmd.RepoName, err)
			return domain.NewError(http.StatusInternalServerError, err)
		}
	}

	sub := domain.NewSubscription().
		WithUser(user).
		WithGitRepository(gitRepo).
		WithNewID()

	subscribeDOIToken := domain.NewDOISubscriptionToken(
		domain.DOISubscriptionTokenActionSubscribe,
	).WithNewID()
	unsubscribeDOIToken := domain.NewDOISubscriptionToken(
		domain.DOISubscriptionTokenActionUnsubscribe,
	).WithNewID()

	sub = sub.
		WithDOISubscriptionToken(subscribeDOIToken).
		WithDOISubscriptionToken(unsubscribeDOIToken)

	confirmActionLink, err := subscribeDOIToken.ToHttpLink(a.cfg.HostURL)
	if err != nil {
		log.Printf("error: failed to generate confirm action link: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	unsubscribeActionLink, err := unsubscribeDOIToken.ToHttpLink(a.cfg.HostURL)
	if err != nil {
		log.Printf("error: failed to generate unsubscribe action link: %v", err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	// sending email and saving subscription should be in the same transaction
	if err := a.repo.WithTransaction(ctx, func(ctx context.Context) error {
		if err := a.repo.SaveRepositorySubscription(ctx, sub, domain.SaveRepositorySubscriptionParams{
			SaveDOITokens: true,
		}); err != nil {
			return err
		}

		if err := a.mailSvc.SendSubscribeRequestEmail(ctx, mailsvc.SubscribeRequestParams{
			Email:                 cmd.Email,
			Repo:                  cmd.RepoName,
			ConfirmActionLink:     confirmActionLink,
			UnsubscribeActionLink: unsubscribeActionLink,
			ConfirmationToken:     subscribeDOIToken.ID().String(),
			UnsubscribeToken:      unsubscribeDOIToken.ID().String(),
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Printf("error: failed to save subscription and send email email=%s: %v", cmd.Email, err)
		return domain.NewError(http.StatusInternalServerError, err)
	}

	return nil
}
