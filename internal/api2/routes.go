package api2

import (
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gin-gonic/gin"
)

func SetupRouter(system *actor.ActorSystem, enginePID *actor.PID) *gin.Engine {
	router := gin.Default()
	// router := gin.New()
	// router.Use(gin.Logger(), gin.Recovery()) // Attach middleware explicitly
	// gin.SetMode(gin.ReleaseMode)

	router.POST("/api/register", func(c *gin.Context) {
		RegisterUserHandler(c, system, enginePID)
	})
	router.GET("/api/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "User registration initiated"})
	})
	router.POST("/api/posts", func(c *gin.Context) {
		CreatePostHandler(c, system, enginePID)
	})
	router.GET("/api/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Post creation initiated"})
	})
	router.POST("/api/newsubreddit", func(c *gin.Context) {
		CreateSubredditHandler(c, system, enginePID)
	})
	router.POST("/api/comment", func(c *gin.Context) {
		CreateCommentHandler(c, system, enginePID)
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Reddit Clone API"})
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	return router
}
