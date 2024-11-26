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
	case *proto.VoteMsg:
		log.Printf("Received VoteMsg: %+v", msg)
		e.handleVote(context, msg)
	default:
		log.Printf("Unhandled message type: %+v", msg)
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

func (e *RedditEngine) handleVote(context actor.Context, msg *proto.VoteMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the vote target is a post
	if post, exists := e.posts[msg.TargetId]; exists {
		applyVoteToPost(post, msg.IsUpvote)
		log.Printf("Vote applied to post: PostID=%s, Upvotes=%d, Downvotes=%d, UserID=%s",
			post.ID, post.Upvotes, post.Downvotes, msg.UserId)
		return
	}

	// Check if the vote target is a comment (very inefficient)
	for _, comments := range e.comments {
		for _, comment := range comments {
			if comment.ID == msg.TargetId {
				applyVoteToComment(comment, msg.IsUpvote)
				log.Printf("Vote applied to comment: CommentID=%s, Upvotes=%d, Downvotes=%d, UserID=%s",
					comment.ID, comment.Upvotes, comment.Downvotes, msg.UserId)
				return
			}
		}
	}

	log.Printf("Target not found for VoteMsg: %+v", msg)
}

func (e *RedditEngine) handleComment(context actor.Context, msg *proto.CreateCommentMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the post exists
	post, postExists := e.posts[msg.PostId]
	if !postExists {
		log.Printf("Post not found for CreateCommentMsg: %+v", msg)
		return
	}

	// Create a new comment
	comment := &Comment{
		ID:        generateID(),
		Content:   msg.Content,
		AuthorID:  msg.AuthorId,
		PostID:    msg.PostId,
		ParentID:  msg.ParentId,
		Children:  []*Comment{},
		CreatedAt: time.Now(),
	}

	// Add to the comments map
	if _, exists := e.comments[comment.PostID]; !exists {
		e.comments[comment.PostID] = []*Comment{}
	}
	e.comments[comment.PostID] = append(e.comments[comment.PostID], comment)

	// If the comment has a parent, link it to the parent
	if comment.ParentID != "" {
		parentFound := false
		// Get all the comments for the post
		for _, c := range e.comments[comment.PostID] {
			// Check if the parent comment exists
			if c.ID == comment.ParentID {
				c.Children = append(c.Children, comment)
				parentFound = true
				break
			}
		}
		if !parentFound {
			log.Printf("Parent comment not found for CommentID=%s, ParentID=%s", comment.ID, comment.ParentID)
			return
		}
	} else {
		// Root-level comment
		post.Comments = append(post.Comments, comment)
	}

	log.Printf("Comment added: %+v", comment)
}

func applyVoteToPost(post *Post, isUpvote bool) {
	if isUpvote {
		post.Upvotes++
	} else {
		post.Downvotes++
	}
}

func applyVoteToComment(comment *Comment, isUpvote bool) {
	if isUpvote {
		comment.Upvotes++
	} else {
		comment.Downvotes++
	}
}

func generateID() string {
	return time.Now().Format("20060102150405.000")
}
