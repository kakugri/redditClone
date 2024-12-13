package api2

import (
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gin-gonic/gin"
	"github.com/kakugri/redditClone/internal/proto"
)

func RegisterUserHandler(c *gin.Context, system *actor.ActorSystem, enginePID *actor.PID) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	system.Root.Send(enginePID, &proto.RegisterUserMsg{Username: req.Username})
	c.JSON(http.StatusOK, gin.H{"message": "User registration initiated"})
}

func CreatePostHandler(c *gin.Context, system *actor.ActorSystem, enginePID *actor.PID) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &proto.CreatePostMsg{
		Title:   req.Title,
		Content: req.Content,
	}
	system.Root.Send(enginePID, msg)
	c.JSON(http.StatusOK, gin.H{"message": "Post creation initiated"})
}

func CreateSubredditHandler(c *gin.Context, system *actor.ActorSystem, enginePID *actor.PID) {
	var req struct {
		Name        string `json:"username"`
		Description string `json:"description"`
		CreatorId   string `json:"creatorid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &proto.CreateSubredditMsg{
		Name:        req.Name,
		Description: req.Description,
		CreatorId:   req.CreatorId,
	}

	system.Root.Send(enginePID, msg)
	c.JSON(http.StatusOK, gin.H{"message": "Subreddit created"})
}

func CreateCommentHandler(c *gin.Context, system *actor.ActorSystem, enginePID *actor.PID) {
	var req struct {
		Content string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &proto.CreateCommentMsg{
		PostId:   "Simulated Post Title",
		Content:  req.Content,
		AuthorId: "Simulated AuthorId",
		ParentId: "Simulated ParentId",
	}

	system.Root.Send(enginePID, msg)
	c.JSON(http.StatusOK, gin.H{"message": "Comment created"})
}
