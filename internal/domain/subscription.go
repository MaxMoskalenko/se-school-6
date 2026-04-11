package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	id             uuid.UUID
	confirmedAt    *time.Time
	unsubscribedAt *time.Time

	user                  *User
	gitRepository         *GitRepository
	doiSubscriptionTokens []*DOISubscriptionToken
}

func NewSubscription() *Subscription {
	return &Subscription{
		doiSubscriptionTokens: make([]*DOISubscriptionToken, 0),
	}
}

func (s *Subscription) WithID(id uuid.UUID) *Subscription {
	s.id = id
	return s
}

func (s *Subscription) WithNewID() *Subscription {
	s.id = uuid.New()
	return s
}

func (s *Subscription) WithConfirmedAt(confirmedAt *time.Time) *Subscription {
	s.confirmedAt = confirmedAt
	return s
}

func (s *Subscription) WithUnsubscribedAt(unsubscribedAt *time.Time) *Subscription {
	s.unsubscribedAt = unsubscribedAt
	return s
}

func (s *Subscription) WithUser(user *User) *Subscription {
	s.user = user
	return s
}

func (s *Subscription) WithGitRepository(repo *GitRepository) *Subscription {
	s.gitRepository = repo
	return s
}

func (s *Subscription) WithDOISubscriptionToken(token *DOISubscriptionToken) *Subscription {
	s.doiSubscriptionTokens = append(s.doiSubscriptionTokens, token)
	return s
}

func (s Subscription) ID() uuid.UUID {
	return s.id
}

func (s Subscription) ConfirmedAt() *time.Time {
	return s.confirmedAt
}

func (s Subscription) UnsubscribedAt() *time.Time {
	return s.unsubscribedAt
}

func (s Subscription) User() *User {
	return s.user
}

func (s Subscription) GitRepository() *GitRepository {
	return s.gitRepository
}

func (s Subscription) DOISubscriptionTokens() []*DOISubscriptionToken {
	return s.doiSubscriptionTokens
}

func (s Subscription) IsActive() bool {
	return s.confirmedAt != nil && s.unsubscribedAt == nil
}

func (s Subscription) IsConfirmed() bool {
	return s.confirmedAt != nil
}

func (s Subscription) DOITokens() []DOISubscriptionTokenAction {
	actions := make([]DOISubscriptionTokenAction, len(s.doiSubscriptionTokens))
	for i, token := range s.doiSubscriptionTokens {
		actions[i] = token.Action()
	}
	return actions
}

func (s *Subscription) WithNewTokens() *Subscription {
	s.doiSubscriptionTokens = []*DOISubscriptionToken{
		NewDOISubscriptionToken(DOISubscriptionTokenActionSubscribe).WithNewID(),
		NewDOISubscriptionToken(DOISubscriptionTokenActionUnsubscribe).WithNewID(),
	}
	return s
}

func (s Subscription) SubscribeToken() *DOISubscriptionToken {
	for _, t := range s.doiSubscriptionTokens {
		if t.Action() == DOISubscriptionTokenActionSubscribe {
			return t
		}
	}
	return nil
}

func (s Subscription) UnsubscribeToken() *DOISubscriptionToken {
	for _, t := range s.doiSubscriptionTokens {
		if t.Action() == DOISubscriptionTokenActionUnsubscribe {
			return t
		}
	}
	return nil
}

func (s Subscription) DOITokenAction(token string) DOISubscriptionTokenAction {
	for _, t := range s.doiSubscriptionTokens {
		if t.ID().String() == token {
			return t.Action()
		}
	}
	return ""
}
