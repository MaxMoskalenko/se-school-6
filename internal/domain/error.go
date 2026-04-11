package domain

import "errors"

var (
	ErrInternal          = errors.New("internal server error")
	ErrNotFound          = errors.New("not found")
	ErrAlreadySubscribed = errors.New("email already subscribed to this repository")
	ErrInvalidToken      = errors.New("invalid token")
	ErrAlreadyConfirmed  = errors.New("subscription already confirmed")
	ErrNotActive         = errors.New("subscription is not active")
	ErrRepoNotFound      = errors.New("repository not found on github")
)

type Error struct {
	code int
	error
}

func NewError(code int, err error) *Error {
	return &Error{
		code:  code,
		error: err,
	}
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.Error()
}
