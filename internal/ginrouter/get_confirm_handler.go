package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type getConfirmRequest struct {
	Token string `uri:"token" binding:"required"`
}

func getConfirmHandler(app *api.App, _ Config, c *gin.Context) {
	var req getConfirmRequest
	if err := c.BindUri(&req); err != nil {
		c.AsciiJSON(400, gin.H{"error": validationErrorMessage(err)})
		return
	}

	if dErr := app.ConfirmSubscription(c.Request.Context(), api.ConfirmSubscriptionCommand{
		Token: req.Token,
	}); dErr != nil {
		c.AsciiJSON(dErr.Code(), gin.H{"error": dErr.Message()})
		return
	}

	c.AsciiJSON(200, gin.H{"status": "confirmed"})
}
