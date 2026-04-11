package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/internal/mockrepo"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestApp(repo *mockrepo.MockRepository, mail *mailsvc.Mock, git *gitsvc.Mock) *App {
	return NewApp(repo, Config{HostURL: "http://localhost:8080"}, mail, git)
}

func TestSubscribeOnRepo_Success(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()
	gitRepo := domain.NewGitRepository("owner", "repo").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(gitRepo, nil)
	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
	repo.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mail.On("SendSubscribeRequestEmail", mock.Anything, mock.Anything).Return(nil)

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.Nil(t, dErr)
	repo.AssertExpectations(t)
	mail.AssertExpectations(t)
}

func TestSubscribeOnRepo_ReadUserFails(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
	assert.Equal(t, domain.ErrInternal.Error(), dErr.Message())
}

func TestSubscribeOnRepo_RepoNotInDB_NotOnGithub(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
	git.On("RepoExists", mock.Anything, "owner", "repo").Return(false, nil)

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusNotFound, dErr.Code())
}

func TestSubscribeOnRepo_RepoNotInDB_GithubCheckFails(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
	git.On("RepoExists", mock.Anything, "owner", "repo").Return(false, fmt.Errorf("github error"))

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}

func TestSubscribeOnRepo_RepoNotInDB_ExistsOnGithub_CreateSucceeds(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()
	gitRepo := domain.NewGitRepository("owner", "repo").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound).Once()
	git.On("RepoExists", mock.Anything, "owner", "repo").Return(true, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(gitRepo, nil).Once()
	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
	repo.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mail.On("SendSubscribeRequestEmail", mock.Anything, mock.Anything).Return(nil)

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.Nil(t, dErr)
	git.AssertCalled(t, "RepoExists", mock.Anything, "owner", "repo")
}

func TestSubscribeOnRepo_AlreadySubscribed(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()
	gitRepo := domain.NewGitRepository("owner", "repo").WithNewID()
	existingSub := domain.NewSubscription().WithNewID().WithUser(user).WithGitRepository(gitRepo)

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(gitRepo, nil)
	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(existingSub, nil)

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusConflict, dErr.Code())
	repo.AssertNotCalled(t, "WithTransaction", mock.Anything, mock.Anything)
}

func TestSubscribeOnRepo_SendEmailFails(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()
	gitRepo := domain.NewGitRepository("owner", "repo").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(gitRepo, nil)
	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)
	repo.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	repo.On("SaveRepositorySubscription", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mail.On("SendSubscribeRequestEmail", mock.Anything, mock.Anything).Return(fmt.Errorf("email error"))

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}

func TestSubscribeOnRepo_ReadSubscriptionFails(t *testing.T) {
	repo := mockrepo.New()
	mail := mailsvc.NewMock()
	git := gitsvc.NewMock()
	app := newTestApp(repo, mail, git)

	user := domain.NewUser("test@example.com").WithNewID()
	gitRepo := domain.NewGitRepository("owner", "repo").WithNewID()

	repo.On("ReadUser", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("ReadGitRepository", mock.Anything, mock.Anything).Return(gitRepo, nil)
	repo.On("ReadRepositorySubscription", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	dErr := app.SubscribeOnRepo(context.Background(), SubscribeOnRepoCommand{
		Email:     "test@example.com",
		RepoOwner: "owner",
		RepoName:  "repo",
	})

	assert.NotNil(t, dErr)
	assert.Equal(t, http.StatusInternalServerError, dErr.Code())
}
