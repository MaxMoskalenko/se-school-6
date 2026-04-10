package mailsvc

import (
	"context"

	"github.com/mrz1836/postmark"
)

type Postmark struct {
	client *postmark.Client
	cfg    PostmarkConfig
}

func NewPostmark(cfg PostmarkConfig) *Postmark {

	client := postmark.NewClient(cfg.ServerToken, cfg.AccountToken)

	return &Postmark{
		client: client,
		cfg:    cfg,
	}
}

func (s *Postmark) SendSubscribeRequestEmail(ctx context.Context, params SubscribeRequestParams) error {
	if _, err := s.client.SendTemplatedEmail(ctx, postmark.TemplatedEmail{
		TemplateID: s.cfg.SubscribeRequestTemplateID,
		TemplateModel: map[string]interface{}{
			"email":                  params.Email,
			"repo":                   params.Repo,
			"confirm_action_url":     params.ConfirmActionLink,
			"unsubscribe_action_url": params.UnsubscribeActionLink,
			"confirmation_token":     params.ConfirmationToken,
			"unsubscribe_token":      params.UnsubscribeToken,
		},
		To:   params.Email,
		From: s.cfg.SenderEmail,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Postmark) SendNewReleaseEmail(ctx context.Context, params NewReleaseEmailParams) error {
	if _, err := s.client.SendTemplatedEmail(ctx, postmark.TemplatedEmail{
		TemplateID: s.cfg.NewReleaseTemplateID,
		TemplateModel: map[string]interface{}{
			"email": params.Email,
			"repo":  params.Repo,
			"tag":   params.ReleaseTag,
		},
		To:   params.Email,
		From: s.cfg.SenderEmail,
	}); err != nil {
		return err
	}

	return nil
}
