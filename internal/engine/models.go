// internal/engine/models.go
package engine

import (
	"time"
)

type User struct {
	ID       string
	Username string
	Karma    int
	JoinDate time.Time
}

type Subreddit struct {
	ID          string
	Name        string
	Description string
	Members     map[string]*User
	Posts       []*Post
	CreatedAt   time.Time
}

type Post struct {
	ID          string
	Title       string
	Content     string
	AuthorID    string
	SubredditID string
	Upvotes     int
	Downvotes   int
	Comments    []*Comment
	CreatedAt   time.Time
}

type Comment struct {
	ID        string
	Content   string
	AuthorID  string
	PostID    string
	ParentID  string // For hierarchical comments
	Children  []*Comment
	Upvotes   int
	Downvotes int
	CreatedAt time.Time
}

type DirectMessage struct {
	ID         string
	FromUserID string
	ToUserID   string
	Content    string
	CreatedAt  time.Time
}
