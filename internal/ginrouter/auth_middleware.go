package ginrouter

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const emailContextKey contextKey = "email"

func authMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AsciiJSON(401, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AsciiJSON(401, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AsciiJSON(401, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			c.AsciiJSON(401, gin.H{"error": "missing email in token"})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), emailContextKey, email)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func EmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(emailContextKey).(string)
	return email, ok
}

func validateAuthEmail(c *gin.Context, paramEmail string, validate bool) bool {
	if !validate {
		return true
	}

	jwtEmail, ok := EmailFromContext(c.Request.Context())
	if !ok || jwtEmail != paramEmail {
		c.AsciiJSON(403, gin.H{"error": "token email does not match request email"})
		c.Abort()
		return false
	}
	return true
}
