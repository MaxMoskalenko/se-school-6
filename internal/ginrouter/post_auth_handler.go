package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type postAuthRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func postAuthHandler(app *api.App, _ Config, c *gin.Context) {
	var req postAuthRequest
	if err := c.BindJSON(&req); err != nil {
		c.AsciiJSON(400, gin.H{"error": validationErrorMessage(err)})
		return
	}

	result, dErr := app.CreateAuthJWT(c.Request.Context(), api.CreateAuthJWTCommand{
		Email: req.Email,
	})
	if dErr != nil {
		c.AsciiJSON(dErr.Code(), gin.H{"error": dErr.Message()})
		return
	}

	c.AsciiJSON(200, gin.H{"token": result.Token})
}
