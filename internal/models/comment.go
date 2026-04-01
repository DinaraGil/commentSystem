package models

import "time"

type Comment struct {
	ID             int       `json:"id" db:"comment_id"`
	ReplyCommentID *int      `json:"replyCommentId,omitempty" db:"reply_comment_id"`
	CommentLevel   int       `json:"commentLevel" db:"comment_level"`
	PostID         int       `json:"postID" db:"post_id"`
	UserID         int       `json:"userID" db:"user_id"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`

	Replies []*Comment `json:"replies"`
}
