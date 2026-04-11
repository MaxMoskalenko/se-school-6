package api

import "errors"

var (
	errInvalidToken     = errors.New("invalid token")
	errAlreadyConfirmed = errors.New("subscription already confirmed")
	errNotActive        = errors.New("subscription is not active")
	errRepoNotFound     = errors.New("repository not found on github")
)
