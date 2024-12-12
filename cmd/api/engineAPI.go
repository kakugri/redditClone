// cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/gin-gonic/gin"
	"github.com/kakugri/redditClone/internal/engine"
	"github.com/kakugri/redditClone/internal/proto"
)

// Global variables for the actor system and the Reddit engine PID
var system *actor.ActorSystem
var enginePID *actor.PID

func main() {
	// Initialize the actor system
	system = actor.NewActorSystem()

	// Configure the remote actor system
	config := remote.Configure("localhost", 8081)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Define the properties for the Reddit engine actor
	props := actor.PropsFromProducer(func() actor.Actor {
		return engine.NewRedditEngine()
	})

	// Spawn the Reddit engine actor and get its PID
	var err error
	enginePID, err = system.Root.SpawnNamed(props, "reddit-engine")
	if err != nil {
		log.Fatalf("Failed to start engine: %v", err)
	}
	log.Printf("RedditEngine started: Address=%s, Id=%s", enginePID.GetAddress(), enginePID.GetId())

	// Initialize the Gin router
	r := gin.Default()

	// Define the routes for the API (Returns JSON response)
	r.POST("/register", registerUser)
	r.POST("/subreddit", createSubreddit)
	r.POST("/post", createPost)
	r.POST("/comment", createComment)
	r.POST("/vote", vote)
	r.POST("/message", sendMessage)

	// Start the Gin server on port 8080
	r.Run(":9090")
}

// Handler for registering a user
func registerUser(c *gin.Context) {
	var msg proto.RegisterUserMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "user registered"})
}

// Handler for creating a subreddit
func createSubreddit(c *gin.Context) {
	var msg proto.CreateSubredditMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "subreddit created"})
}

// Handler for creating a post
func createPost(c *gin.Context) {
	var msg proto.CreatePostMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "post created"})
}

// Handler for creating a comment
func createComment(c *gin.Context) {
	var msg proto.CreateCommentMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "comment created"})
}

// Handler for voting on a post or comment
func vote(c *gin.Context) {
	var msg proto.VoteMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "vote recorded"})
}

// Handler for sending a direct message
func sendMessage(c *gin.Context) {
	var msg proto.DirectMessageMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	system.Root.Send(enginePID, &msg)
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}
