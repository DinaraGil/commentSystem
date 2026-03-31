package graph

import (
	"commentSystem/internal/models"
	"context"
	"database/sql"
	"fmt"
)

func (r *Resolver) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	query := `SELECT person_id, username FROM person WHERE person_id=$1`
	err := r.DB.Get(&user, query, userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with id %d not found", userID)
	}
	return &user, err
}

func (r *Resolver) GetPostByID(ctx context.Context, postId int) (*models.Post, error) {
	var post models.Post

	query := `SELECT * FROM post WHERE post_id=$1`
	err := r.DB.Get(&post, query, postId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post with id %d not found", postId)
	}
	return &post, err
}
