package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type getUnsubscribeRequest struct {
	Token string `uri:"token" binding:"required"`
}

func getUnsubscribeHandler(app *api.App, c *gin.Context) {
	var req getUnsubscribeRequest
	if err := c.BindUri(&req); err != nil {
		c.AsciiJSON(400, gin.H{"error": validationErrorMessage(err)})
		return
	}

	if dErr := app.UnsubscribeFromRepo(c.Request.Context(), api.UnsubscribeCommand{
		Token: req.Token,
	}); dErr != nil {
		c.AsciiJSON(dErr.Code(), gin.H{"error": dErr.Message()})
		return
	}

	c.AsciiJSON(200, gin.H{"status": "unsubscribed"})
}
