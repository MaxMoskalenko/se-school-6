package scanner

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/MaxMoskalenko/se-school-6/internal/mockrepo"
	"github.com/MaxMoskalenko/se-school-6/pkg/gitsvc"
	"github.com/MaxMoskalenko/se-school-6/pkg/mailsvc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestApp(repo *mockrepo.MockRepository, git *gitsvc.Mock, mail *mailsvc.Mock) *App {
	return NewApp(repo, Config{Interval: time.Minute}, git, mail)
}

func buildRepo(owner, name string, lastSeenTag *string, subs ...*domain.Subscription) *domain.GitRepository {
	repo := domain.NewGitRepository(owner, name).WithNewID()
	if lastSeenTag != nil {
		repo = repo.WithLastSeenTag(*lastSeenTag)
	}
	for _, sub := range subs {
		repo = repo.AttachSubscription(sub)
	}
	return repo
}

func buildSubWithUser(email string) *domain.Subscription {
	user := domain.NewUser(email).WithNewID()
	return domain.NewSubscription().WithNewID().WithUser(user)
}

func TestScan_NoRepos(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{}, nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	git.AssertNotCalled(t, "FetchLatestReleaseTag")
}

func TestScan_ReadReposFails(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("db error"))

	err := app.Scan(context.Background())

	assert.Error(t, err)
}

func TestScan_NewRelease_NotifiesSubscribers(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	oldTag := "v1.0.0"
	sub := buildSubWithUser("user@example.com")
	gitRepo := buildRepo("owner", "repo", &oldTag, sub)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{gitRepo}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo").Return("v2.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.MatchedBy(func(p mailsvc.NewReleaseEmailParams) bool {
		return p.Email == "user@example.com" && p.Repo == "owner/repo" && p.ReleaseTag == "v2.0.0"
	})).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	mail.AssertExpectations(t)
	repo.AssertCalled(t, "SaveGitRepository", mock.Anything, mock.Anything)
}

func TestScan_NoNewRelease_NoNotification(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag := "v1.0.0"
	sub := buildSubWithUser("user@example.com")
	gitRepo := buildRepo("owner", "repo", &tag, sub)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{gitRepo}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo").Return("v1.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	mail.AssertNotCalled(t, "SendNewReleaseEmail")
	repo.AssertCalled(t, "SaveGitRepository", mock.Anything, mock.Anything)
}

func TestScan_FirstRelease_NilLastSeenTag(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	sub := buildSubWithUser("user@example.com")
	gitRepo := buildRepo("owner", "repo", nil, sub)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{gitRepo}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo").Return("v1.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.Anything).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	mail.AssertCalled(t, "SendNewReleaseEmail", mock.Anything, mock.Anything)
}

func TestScan_RateLimited_AbortsImmediately(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag := "v1.0.0"
	repo1 := buildRepo("owner", "repo1", &tag)
	repo2 := buildRepo("owner", "repo2", &tag)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{repo1, repo2}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo1").Return("", gitsvc.ErrRateLimited)

	err := app.Scan(context.Background())

	assert.ErrorIs(t, err, gitsvc.ErrRateLimited)
	// repo2 should never be checked
	git.AssertNotCalled(t, "FetchLatestReleaseTag", mock.Anything, "owner", "repo2")
}

func TestScan_FetchTagFails_ContinuesToNextRepo(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag := "v1.0.0"
	sub := buildSubWithUser("user@example.com")
	repo1 := buildRepo("owner", "repo1", &tag)
	repo2 := buildRepo("owner", "repo2", &tag, sub)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{repo1, repo2}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo1").Return("", fmt.Errorf("network error"))
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo2").Return("v2.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.Anything).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	mail.AssertCalled(t, "SendNewReleaseEmail", mock.Anything, mock.Anything)
}

func TestScan_SaveRepoFails_ContinuesToNextRepo(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag := "v1.0.0"
	sub := buildSubWithUser("user@example.com")
	repo1 := buildRepo("owner", "repo1", &tag)
	repo2 := buildRepo("owner", "repo2", &tag, sub)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{repo1, repo2}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo1").Return("v2.0.0", nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo2").Return("v2.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.MatchedBy(func(r *domain.GitRepository) bool {
		return r.Name() == "repo1"
	})).Return(fmt.Errorf("save error"))
	repo.On("SaveGitRepository", mock.Anything, mock.MatchedBy(func(r *domain.GitRepository) bool {
		return r.Name() == "repo2"
	})).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.Anything).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	// repo1 save failed, so no email for it; repo2 succeeds
	mail.AssertNumberOfCalls(t, "SendNewReleaseEmail", 1)
}

func TestScan_SendEmailFails_ContinuesToNextSubscriber(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag := "v1.0.0"
	sub1 := buildSubWithUser("user1@example.com")
	sub2 := buildSubWithUser("user2@example.com")
	gitRepo := buildRepo("owner", "repo", &tag, sub1, sub2)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{gitRepo}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo").Return("v2.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.MatchedBy(func(p mailsvc.NewReleaseEmailParams) bool {
		return p.Email == "user1@example.com"
	})).Return(fmt.Errorf("email error"))
	mail.On("SendNewReleaseEmail", mock.Anything, mock.MatchedBy(func(p mailsvc.NewReleaseEmailParams) bool {
		return p.Email == "user2@example.com"
	})).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	mail.AssertNumberOfCalls(t, "SendNewReleaseEmail", 2)
}

func TestScan_MultipleRepos_MultipleSubscribers(t *testing.T) {
	repo := mockrepo.New()
	git := gitsvc.NewMock()
	mail := mailsvc.NewMock()
	app := newTestApp(repo, git, mail)

	tag1 := "v1.0.0"
	tag2 := "v3.0.0"
	sub1 := buildSubWithUser("alice@example.com")
	sub2 := buildSubWithUser("bob@example.com")
	// repo1 has a new release, repo2 does not
	gitRepo1 := buildRepo("owner", "repo1", &tag1, sub1, sub2)
	gitRepo2 := buildRepo("owner", "repo2", &tag2, sub1)

	repo.On("ReadGitRepositories", mock.Anything, mock.Anything).Return([]*domain.GitRepository{gitRepo1, gitRepo2}, nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo1").Return("v2.0.0", nil)
	git.On("FetchLatestReleaseTag", mock.Anything, "owner", "repo2").Return("v3.0.0", nil)
	repo.On("SaveGitRepository", mock.Anything, mock.Anything).Return(nil)
	mail.On("SendNewReleaseEmail", mock.Anything, mock.Anything).Return(nil)

	err := app.Scan(context.Background())

	assert.NoError(t, err)
	// 2 emails for repo1 (alice + bob), 0 for repo2 (no new release)
	mail.AssertNumberOfCalls(t, "SendNewReleaseEmail", 2)
}
