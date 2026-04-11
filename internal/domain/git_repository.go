package domain

import (
	"time"

	"github.com/google/uuid"
)

type GitRepository struct {
	id            uuid.UUID
	owner         string
	name          string
	lastSeenTag   *string
	lastCheckedAt *time.Time
	subscriptions []*Subscription
}

func NewGitRepository(owner, name string) *GitRepository {
	return &GitRepository{
		owner: owner,
		name:  name,
	}
}

func (r *GitRepository) WithID(id uuid.UUID) *GitRepository {
	r.id = id
	return r
}

func (r *GitRepository) WithNewID() *GitRepository {
	r.id = uuid.New()
	return r
}

func (r GitRepository) ID() uuid.UUID {
	return r.id
}

func (r *GitRepository) WithLastSeenTag(tag string) *GitRepository {
	r.lastSeenTag = &tag
	return r
}

func (r GitRepository) Owner() string {
	return r.owner
}

func (r GitRepository) Name() string {
	return r.name
}

func (r GitRepository) LastSeenTag() *string {
	return r.lastSeenTag
}

func (r *GitRepository) WithLastCheckedAt(t *time.Time) *GitRepository {
	r.lastCheckedAt = t
	return r
}

func (r GitRepository) LastCheckedAt() *time.Time {
	return r.lastCheckedAt
}

func (r *GitRepository) AttachSubscription(sub *Subscription) *GitRepository {
	r.subscriptions = append(r.subscriptions, sub)
	return r
}

func (r GitRepository) Subscriptions() []*Subscription {
	return r.subscriptions
}
