package mailsvc

type PostmarkConfig struct {
	ServerToken  string
	AccountToken string

	SenderEmail string

	SubscribeRequestTemplateID int64
	NewReleaseTemplateID       int64
}

type SubscribeRequestParams struct {
	Email string
	Repo  string

	ConfirmationToken string
	UnsubscribeToken  string

	ConfirmActionLink     string
	UnsubscribeActionLink string
}

type NewReleaseEmailParams struct {
	Email string
	Repo  string

	ReleaseTag string
}
