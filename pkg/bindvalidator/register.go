package bindvalidator

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var repoRegex = regexp.MustCompile(`^[0-9a-zA-Z\-_.]+/[0-9a-zA-Z\-_.]+$`)

func Register() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("failed to get validator engine")
	}

	_ = v.RegisterValidation("repo", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return repoRegex.MatchString(val)
	})

	return nil
}
