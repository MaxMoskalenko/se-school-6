package api

import (
	"context"
	"net/http"
	"time"

	"github.com/MaxMoskalenko/se-school-6/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type CreateAuthJWTCommand struct {
	Email string
}

type CreateAuthJWTResult struct {
	Token string
}

func (a *App) CreateAuthJWT(_ context.Context, cmd CreateAuthJWTCommand) (*CreateAuthJWTResult, *domain.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": cmd.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return nil, domain.NewError(http.StatusInternalServerError, domain.ErrInternal)
	}

	return &CreateAuthJWTResult{Token: signed}, nil
}
