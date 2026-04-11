package domain

import (
	"strings"

	"github.com/google/uuid"
)

// using string for action to make it more readable in the logs
type DOISubscriptionTokenAction string

const (
	DOISubscriptionTokenActionUnknown     DOISubscriptionTokenAction = "unknown"
	DOISubscriptionTokenActionSubscribe   DOISubscriptionTokenAction = "subscribe"
	DOISubscriptionTokenActionUnsubscribe DOISubscriptionTokenAction = "unsubscribe"
)

type DOISubscriptionToken struct {
	id     uuid.UUID
	action DOISubscriptionTokenAction
}

func NewDOISubscriptionToken(action DOISubscriptionTokenAction) *DOISubscriptionToken {
	return &DOISubscriptionToken{
		action: action,
	}
}

func (t *DOISubscriptionToken) WithID(id uuid.UUID) *DOISubscriptionToken {
	t.id = id
	return t
}

func (t *DOISubscriptionToken) WithNewID() *DOISubscriptionToken {
	t.id = uuid.New()
	return t
}

func (t DOISubscriptionToken) ID() uuid.UUID {
	return t.id
}

func (t DOISubscriptionToken) Action() DOISubscriptionTokenAction {
	return t.action
}

func (t DOISubscriptionToken) ToHttpLink(url string) (string, error) {
	trimmed := strings.TrimSuffix(url, "/")

	switch t.action {
	case DOISubscriptionTokenActionSubscribe:
		return trimmed + "/api/confirm/" + t.ID().String(), nil
	case DOISubscriptionTokenActionUnsubscribe:
		return trimmed + "/api/unsubscribe/" + t.ID().String(), nil
	default:
		return "", nil
	}
}
