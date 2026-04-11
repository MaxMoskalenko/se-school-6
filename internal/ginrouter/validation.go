package ginrouter

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var fieldMessages = map[string]map[string]string{
	"Email": {
		"required": "email is required",
		"email":    "invalid email format",
	},
	"Repo": {
		"required": "repo is required",
		"repo":     "invalid repo format, expected 'owner/repo'",
	},
	"Token": {
		"required": "token is required",
	},
}

func validationErrorMessage(err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return "invalid request"
	}

	fe := ve[0]
	if msgs, ok := fieldMessages[fe.Field()]; ok {
		if msg, ok := msgs[fe.Tag()]; ok {
			return msg
		}
	}

	return fmt.Sprintf("invalid value for field '%s'", fe.Field())
}
