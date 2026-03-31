package models

import (
	"time"
)

type Post struct {
	ID           int       `json:"id" db:"post_id"`
	UserId       int       `json:"userId" db:"person_id"`
	Content      string    `json:"content" db:"content"`
	AllowComment *bool     `json:"allowComment,omitempty" db:"allow_comment"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`

	Comments []*Comment `json:"replies"`
}
