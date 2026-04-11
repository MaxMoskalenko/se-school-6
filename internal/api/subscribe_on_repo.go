package api

import (
	"context"
	"errors"
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
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	gitRepo, err := a.repo.ReadGitRepository(ctx, domain.ReadGitRepositoryParams{
		ByOwner: &cmd.RepoOwner,
		ByName:  &cmd.RepoName,
	})
	if err != nil {
		exists, ghErr := a.gitSvc.RepoExists(ctx, cmd.RepoOwner, cmd.RepoName)
		if ghErr != nil {
			log.Printf("error: failed to check repo on github owner=%s name=%s: %v", cmd.RepoOwner, cmd.RepoName, ghErr)
			return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
		}
		if !exists {
			return domain.NewError(http.StatusNotFound, domain.ErrRepoNotFound)
		}

		gitRepoToCreate := domain.NewGitRepository(cmd.RepoOwner, cmd.RepoName).WithNewID()
		gitRepo, err = a.repo.ReadGitRepository(ctx, domain.ReadGitRepositoryParams{
			ByOwner:           &cmd.RepoOwner,
			ByName:            &cmd.RepoName,
			CreateIfNotExists: gitRepoToCreate,
		})
		if err != nil {
			log.Printf("error: failed to create git repository owner=%s name=%s: %v", cmd.RepoOwner, cmd.RepoName, err)
			return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
		}
	}

	userID := user.ID().String()
	repoID := gitRepo.ID().String()

	sub, err := a.repo.ReadRepositorySubscription(ctx, domain.ReadRepositorySubscriptionParams{
		ByUserID:            &userID,
		ByRepositoryID:      &repoID,
		OnlyNonUnsubscribed: true,
	})
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Printf("error: failed to read or create subscription email=%s: %v", cmd.Email, err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	// if subscription already exists and is active, return conflict
	if err == nil {
		return domain.NewError(http.StatusConflict, domain.ErrAlreadySubscribed)
	}

	sub = domain.NewSubscription().
		WithUser(user).
		WithGitRepository(gitRepo).
		WithNewID().
		WithNewTokens()

	confirmActionLink, err := sub.SubscribeToken().ToHttpLink(a.cfg.HostURL)
	if err != nil {
		log.Printf("error: failed to generate confirm action link: %v", err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	unsubscribeActionLink, err := sub.UnsubscribeToken().ToHttpLink(a.cfg.HostURL)
	if err != nil {
		log.Printf("error: failed to generate unsubscribe action link: %v", err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
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
			ConfirmationToken:     sub.SubscribeToken().ID().String(),
			UnsubscribeToken:      sub.UnsubscribeToken().ID().String(),
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Printf("error: failed to save subscription and send email email=%s: %v", cmd.Email, err)
		return domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	return nil
}
