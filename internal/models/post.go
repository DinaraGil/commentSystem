package models

import (
	"time"
)

type Post struct {
	ID           int       `json:"id" db:"post_id"`
	UserID       int       `json:"userID" db:"user_id"`
	Content      string    `json:"content" db:"content"`
	AllowComment *bool     `json:"allowComment,omitempty" db:"allow_comment"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
