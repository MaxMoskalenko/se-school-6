package ginrouter

import (
	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type getSubscriptionsRequest struct {
	Email string `form:"email" binding:"required,email"`
}

type subscriptionResponse struct {
	Email       string  `json:"email"`
	Repo        string  `json:"repo"`
	Confirmed   bool    `json:"confirmed"`
	LastSeenTag *string `json:"last_seen_tag"`
}

func getSubscriptionsHandler(app *api.App, c *gin.Context) {
	var req getSubscriptionsRequest
	if err := c.BindQuery(&req); err != nil {
		c.AsciiJSON(400, gin.H{"error": err.Error()})
		return
	}

	subs, err := app.GetSubscriptions(c.Request.Context(), api.GetSubscriptionsQuery{
		Email: req.Email,
	})
	if err != nil {
		c.AsciiJSON(500, gin.H{"error": err.Error()})
		return
	}

	result := make([]subscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		resp := subscriptionResponse{
			Confirmed: sub.IsConfirmed(),
		}
		if user := sub.User(); user != nil {
			resp.Email = user.Email()
		}
		if repo := sub.GitRepository(); repo != nil {
			resp.Repo = repo.Owner() + "/" + repo.Name()
			resp.LastSeenTag = repo.LastSeenTag()
		}
		result = append(result, resp)
	}

	c.AsciiJSON(200, result)
}
