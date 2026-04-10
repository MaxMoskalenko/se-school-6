package ginrouter

import (
	"fmt"
	"strings"

	"github.com/MaxMoskalenko/se-school-6/internal/api"
	"github.com/gin-gonic/gin"
)

type postSubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Repo  string `json:"repo" binding:"required,repo"`
}

func postSubscribeHandler(app *api.App, c *gin.Context) {
	var req postSubscribeRequest
	if err := c.BindJSON(&req); err != nil {
		c.AsciiJSON(400, gin.H{"error": err.Error()})
		return
	}

	repoOwner, repoName, err := parseRepo(req.Repo)
	if err != nil {
		c.AsciiJSON(400, gin.H{"error": err.Error()})
		return
	}

	if dErr := app.SubscribeOnRepo(c.Request.Context(), api.SubscribeOnRepoCommand{
		Email:     req.Email,
		RepoOwner: repoOwner,
		RepoName:  repoName,
	}); dErr != nil {
		c.AsciiJSON(dErr.Code(), gin.H{"error": dErr.Message()})
		return
	}
	c.AsciiJSON(200, gin.H{"status": "subscribed"})
}

func parseRepo(repo string) (owner, name string, err error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo format, expected 'owner/name'")
	}
	return parts[0], parts[1], nil
}
