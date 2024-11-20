// internal/engine/engine.go
package engine

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type RedditEngine struct {
	users      map[string]*User
	subreddits map[string]*Subreddit
	messages   map[string][]*DirectMessage
	mu         sync.RWMutex
}

func NewRedditEngine() *RedditEngine {
	return &RedditEngine{
		users:      make(map[string]*User),
		subreddits: make(map[string]*Subreddit),
		messages:   make(map[string][]*DirectMessage),
	}
}

func (e *RedditEngine) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *RegisterUserMsg:
		e.handleRegisterUser(context, msg)
	case *CreateSubredditMsg:
		e.handleCreateSubreddit(context, msg)
		// case *CreatePostMsg:
		// 	e.handleCreatePost(context, msg)
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
	context.Respond(user)
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
	context.Respond(subreddit)
}

// Add other handler methods...

func generateID() string {
	return time.Now().String()
}
