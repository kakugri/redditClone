// internal/engine/messages.go
package engine

// Actor Messages
type RegisterUserMsg struct {
	Username string
}

type CreateSubredditMsg struct {
	Name        string
	Description string
	CreatorID   string
}

type JoinSubredditMsg struct {
	UserID      string
	SubredditID string
}

type CreatePostMsg struct {
	Title       string
	Content     string
	AuthorID    string
	SubredditID string
}

type CreateCommentMsg struct {
	Content  string
	AuthorID string
	PostID   string
	ParentID string
}

type VoteMsg struct {
	UserID   string
	TargetID string
	IsUpvote bool
}

type GetFeedMsg struct {
	UserID string
}

type DirectMessageMsg struct {
	FromUserID string
	ToUserID   string
	Content    string
}
