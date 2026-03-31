package models

type NewPost struct {
	UserID       int    `json:"userId" db:"person_id"`
	Content      string `json:"content" db:"content"`
	AllowComment *bool  `json:"allowComment,omitempty" db:"allow_comment"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
}
