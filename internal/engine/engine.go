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
	metrics    *Metrics
	mu         sync.RWMutex
}

func NewRedditEngine() *RedditEngine {
	return &RedditEngine{
		users:      make(map[string]*User),
		subreddits: make(map[string]*Subreddit),
		posts:      make(map[string]*Post),
		messages:   make(map[string][]*DirectMessage),
		comments:   make(map[string][]*Comment),
		metrics: &Metrics{
			StartTime: time.Now(),
		},
	}
}

func (e *RedditEngine) Receive(context actor.Context) {
	message := context.Message()
	log.Printf("Engine received message of type: %T, content: %+v", message, message)

	switch msg := message.(type) {
	case *actor.Started:
		log.Println("RedditEngine started and ready to receive messages.")
		e.StartMetricsReporter(context)
	case *proto.MetricsReportMsg:
		log.Printf("Metrics Report: TotalPosts=%d, ActiveUsers=%d, TotalVotes=%d, TotalComments=%d, TotalMessages=%d", e.metrics.TotalPosts, e.metrics.ActiveUsers,
			e.metrics.TotalVotes, e.metrics.TotalComments, e.metrics.TotalMessages)
	case *proto.RegisterUserMsg:
		log.Printf("Received RegisterUserMsg: %+v", msg)
		e.handleRegisterUser(context, msg)
	case *proto.CreateSubredditMsg:
		log.Printf("Received CreateSubredditMsg: %+v", msg)
		e.handleCreateSubreddit(context, msg)
	case *proto.CreatePostMsg:
		log.Printf("Received CreatePostMsg: %+v", msg)
		e.handleCreatePost(context, msg)
	case *proto.CreateCommentMsg:
		log.Printf("Received CreateCommentMsg: %+v", msg)
		e.handleCreateComment(context, msg)
	case *proto.VoteMsg:
		log.Printf("Received VoteMsg: %+v", msg)
		e.handleVote(context, msg)
	case *proto.DirectMessageMsg:
		log.Printf("Received DirectMessageMsg: %+v", msg)
		e.handleDirectMessage(context, msg)
	default:
		log.Printf("Unhandled message type: %+v", msg)
	}
}

func (e *RedditEngine) StartMetricsReporter(context actor.Context) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			e.ReportMetrics(context)
		}
	}()
}

func (e *RedditEngine) ReportMetrics(context actor.Context) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	metricsCopy := *e.metrics // Create a copy to avoid race conditions
	context.Send(context.Self(), &proto.MetricsReportMsg{
		TotalPosts:    metricsCopy.TotalPosts,
		TotalComments: metricsCopy.TotalComments,
		TotalVotes:    metricsCopy.TotalVotes,
		ActiveUsers:   metricsCopy.ActiveUsers,
		TotalMessages: metricsCopy.TotalMessages,
	})
}

func (e *RedditEngine) updateMetrics(metricFunc func(*Metrics)) {
	e.metrics.mu.Lock()
	defer e.metrics.mu.Unlock()
	metricFunc(e.metrics)
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
	e.updateMetrics(func(m *Metrics) {
		m.ActiveUsers++
	})
	log.Printf("User registered: %+v", user)
}

func (e *RedditEngine) handleCreatePost(context actor.Context, msg *proto.CreatePostMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	post := &Post{
		ID:          generateID(),
		Title:       msg.Title,
		Content:     msg.Content,
		AuthorID:    msg.AuthorId,
		SubredditID: msg.SubredditId,
		CreatedAt:   time.Now(),
	}
	e.posts[post.ID] = post
	e.updateMetrics(func(m *Metrics) {
		m.TotalPosts++
	})
	log.Printf("Post created: %+v", post)
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
}

func (e *RedditEngine) handleDirectMessage(context actor.Context, msg *proto.DirectMessageMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	dm := &DirectMessage{
		ID:         generateID(),
		FromUserID: msg.FromUserId,
		ToUserID:   msg.ToUserId,
		Content:    msg.Content,
		CreatedAt:  time.Now(),
	}
	e.messages[msg.ToUserId] = append(e.messages[msg.ToUserId], dm)
	e.updateMetrics(func(m *Metrics) {
		m.TotalMessages++
	})
	log.Printf("Direct message sent: %+v", dm)
}

func (e *RedditEngine) handleVote(context actor.Context, msg *proto.VoteMsg) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the vote target is a post
	if post, exists := e.posts[msg.TargetId]; exists {
		if msg.IsUpvote {
			post.Upvotes++
		} else {
			post.Downvotes++
		}
		log.Printf("Vote applied to post: PostID=%s, Upvotes=%d, Downvotes=%d, UserID=%s",
			post.ID, post.Upvotes, post.Downvotes, msg.UserId)
		e.updateMetrics(func(m *Metrics) {
			m.TotalVotes++
		})
		return
	}

	// Check if the vote target is a comment (very inefficient)
	for _, comments := range e.comments {
		for _, comment := range comments {
			if comment.ID == msg.TargetId {
				if msg.IsUpvote {
					comment.Upvotes++
				} else {
					comment.Downvotes++
				}
				log.Printf("Vote applied to comment: CommentID=%s, Upvotes=%d, Downvotes=%d, UserID=%s",
					comment.ID, comment.Upvotes, comment.Downvotes, msg.UserId)
				e.updateMetrics(func(m *Metrics) {
					m.TotalVotes++
				})
				return
			}
		}
	}

	log.Printf("Target not found for VoteMsg: %+v", msg)
}

// func (e *RedditEngine) handleVote(context actor.Context, msg *proto.VoteMsg) {
// 	e.mu.Lock()
// 	defer e.mu.Unlock()

// 	post, exists := e.posts[msg.TargetId]
// 	if !exists {
// 		log.Printf("Post not found for vote: %s", msg.TargetId)
// 		return
// 	}

// 	if msg.IsUpvote {
// 		post.Upvotes++
// 	} else {
// 		post.Downvotes++
// 	}
// 	e.updateMetrics(func(m *Metrics) {
// 		m.TotalVotes++
// 	})
// 	log.Printf("Vote processed for post %s: Upvotes=%d, Downvotes=%d", msg.TargetId, post.Upvotes, post.Downvotes)
// }

func (e *RedditEngine) handleCreateComment(context actor.Context, msg *proto.CreateCommentMsg) {
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
	e.updateMetrics(func(m *Metrics) {
		m.TotalComments++
	})
}

// func (e *RedditEngine) handleCreateComment(context actor.Context, msg *proto.CreateCommentMsg) {
// 	e.mu.Lock()
// 	defer e.mu.Unlock()

// 	comment := &Comment{
// 		ID:        generateID(),
// 		Content:   msg.Content,
// 		AuthorID:  msg.AuthorId,
// 		PostID:    msg.PostId,
// 		ParentID:  msg.ParentId,
// 		CreatedAt: time.Now(),
// 	}
// 	e.comments[msg.PostId] = append(e.comments[msg.PostId], comment)
// 	e.updateMetrics(func(m *Metrics) {
// 		m.TotalComments++
// 	})
// 	log.Printf("Comment created: %+v", comment)
// }

func generateID() string {
	return time.Now().Format("20060102150405.000")
}
