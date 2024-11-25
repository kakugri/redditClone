// internal/engine/engine.go
package engine

import (
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type RedditEngine struct {
	users      map[string]*User
	subreddits map[string]*Subreddit
	posts      map[string]*Post
	messages   map[string][]*DirectMessage
	mu         sync.RWMutex
}

func NewRedditEngine() *RedditEngine {
	return &RedditEngine{
		users:      make(map[string]*User),
		subreddits: make(map[string]*Subreddit),
		posts:      make(map[string]*Post),
		messages:   make(map[string][]*DirectMessage),
	}
}

func (e *RedditEngine) Receive(context actor.Context) {
	log.Printf("Engine received message: %+v", context.Message()) // Log all messages
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("RedditEngine started and ready to receive messages.")
	case *RegisterUserMsg:
		log.Printf("Received RegisterUserMsg: %+v", msg)
		e.handleRegisterUser(context, msg)
	case *CreateSubredditMsg:
		log.Printf("Received CreateSubredditMsg: %+v", msg)
		e.handleCreateSubreddit(context, msg)
	case *CreatePostMsg:
		log.Printf("Received CreatePostMsg: %+v", msg)
		e.handleCreatePost(context, msg)
	default:
		log.Printf("Unhandled message type: %+v", msg)
		// case *VoteMsg:
		// 	e.handleVote(context, msg)
		// case *CreateCommentMsg:
		// 	e.handleCreateComment(context, msg)
		// case *DirectMessageMsg:
		// 	e.handleDirectMessage(context, msg)
	}
}

func (e *RedditEngine) handleRegisterUser(context actor.Context, msg *RegisterUserMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user := &User{
		ID:       generateID(),
		Username: msg.Username,
		JoinDate: time.Now(),
	}
	e.users[user.ID] = user
	log.Printf("User registered: %+v", user)
	context.Respond(user)
}

func (e *RedditEngine) handleCreatePost(context actor.Context, msg *CreatePostMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if msg.AuthorID == "" || msg.SubredditID == "" {
		log.Printf("CreatePostMsg missing AuthorID or SubredditID: %+v", msg)
		return
	}

	post := &Post{
		ID:          generateID(),
		Title:       msg.Title,
		Content:     msg.Content,
		AuthorID:    msg.AuthorID,
		SubredditID: msg.SubredditID,
		CreatedAt:   time.Now(),
	}
	e.posts[post.ID] = post
	log.Printf("Post created: %+v", post)
	context.Respond(post)
}

func (e *RedditEngine) handleCreateSubreddit(context actor.Context, msg *CreateSubredditMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	subreddit := &Subreddit{
		ID:          generateID(),
		Name:        msg.Name,
		Description: msg.Description,
		Members:     make(map[string]*User),
		CreatedAt:   time.Now(),
	}
	e.subreddits[subreddit.ID] = subreddit
	log.Printf("Subreddit created: %+v", subreddit)
	context.Respond(subreddit)
}

// Add other handler methods...

func generateID() string {
	return time.Now().Format("20060102150405.000")
}
