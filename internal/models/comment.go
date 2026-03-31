package models

import "time"

type Comment struct {
	ID             int       `json:"id" db:"comment_id"`
	ReplyCommentID *int      `json:"replyCommentId,omitempty" db:"reply_comment_id"`
	PostID         int       `json:"postId" db:"post_id"`
	UserID         int       `json:"userId" db:"person_id"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`

	Replies []*Comment `json:"replies"`
}
