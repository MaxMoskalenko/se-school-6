package mailsvc

import "context"

type Interface interface {
	SendSubscribeRequestEmail(ctx context.Context, params SubscribeRequestParams) error
	SendNewReleaseEmail(ctx context.Context, params NewReleaseEmailParams) error
}
