// internal/engine/engine.go
package engine

import (
	"log"
	"sync"
	"time"

	"github.com/kakugri/redditClone/internal/proto"

	"github.com/asynkron/protoactor-go/actor"
)

type RedditEngine struct {
	users      map[string]*User
	subreddits map[string]*Subreddit
	posts      map[string]*Post
	messages   map[string][]*DirectMessage
	comments   map[string][]*Comment
	mu         sync.RWMutex
}

func NewRedditEngine() *RedditEngine {
	return &RedditEngine{
		users:      make(map[string]*User),
		subreddits: make(map[string]*Subreddit),
		posts:      make(map[string]*Post),
		messages:   make(map[string][]*DirectMessage),
		comments:   make(map[string][]*Comment),
	}
}

func (e *RedditEngine) Receive(context actor.Context) {
	message := context.Message()

	// Log the type and content of the message
	log.Printf("Engine received message of type: %T, content: %+v", message, message)

	switch msg := message.(type) {
	case *actor.Started:
		log.Println("RedditEngine started and ready to receive messages.")
	case *proto.RegisterUserMsg:
		log.Printf("Received RegisterUserMsg: %+v", msg)
		e.handleRegisterUser(context, msg)
	case *proto.CreateSubredditMsg:
		log.Printf("Received CreateSubredditMsg: %+v", msg)
		e.handleCreateSubreddit(context, msg)
	case *proto.CreatePostMsg:
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

func (e *RedditEngine) handleRegisterUser(context actor.Context, msg *proto.RegisterUserMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user := &User{
		ID:       generateID(),
		Username: msg.Username,
		JoinDate: time.Now(),
	}
	e.users[user.ID] = user
	log.Printf("User registered: %+v", user)
	// context.Respond(user)
}

func (e *RedditEngine) handleCreatePost(context actor.Context, msg *proto.CreatePostMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if msg.AuthorId == "" || msg.SubredditId == "" {
		log.Printf("CreatePostMsg missing AuthorID or SubredditID: %+v", msg)
		return
	}

	post := &Post{
		ID:          generateID(),
		Title:       msg.Title,
		Content:     msg.Content,
		AuthorID:    msg.AuthorId,
		SubredditID: msg.SubredditId,
		CreatedAt:   time.Now(),
	}
	e.posts[post.ID] = post
	log.Printf("Post created: %+v", post)
	// context.Respond(post)
}

func (e *RedditEngine) handleCreateSubreddit(context actor.Context, msg *proto.CreateSubredditMsg) {
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
	// context.Respond(subreddit)
}

// Add other handler methods...

func generateID() string {
	return time.Now().Format("20060102150405.000")
}
