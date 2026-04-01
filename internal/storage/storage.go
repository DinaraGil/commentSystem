package storage

import (
	"commentSystem/graph/model"
	"commentSystem/internal/models"
	"context"
)

type Storage interface {
	CreateUser(ctx context.Context, input model.NewUser) (*models.User, error)
	CreatePost(ctx context.Context, input model.NewPost) (*models.Post, error)
	CreateComment(ctx context.Context, input model.NewComment) (*models.Comment, error)

	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetPostByID(ctx context.Context, id int) (*models.Post, error)

	GetPosts(ctx context.Context) ([]*models.Post, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error)
}
