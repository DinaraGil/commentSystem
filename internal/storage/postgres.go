package storage

import (
	"commentSystem/graph/model"
	"commentSystem/internal/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	var user models.User

	query := `INSERT INTO user_data (username) VALUES ($1) RETURNING user_id, username`
	err := s.db.QueryRowxContext(ctx, query, input.Username).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	query := `SELECT user_id, username FROM user_data WHERE user_id = $1`
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (s *PostgresStorage) CreatePost(ctx context.Context, input model.NewPost) (*models.Post, error) {
	if _, err := s.GetUserByID(ctx, input.UserID); err != nil {
		return nil, err
	}

	var post models.Post
	query := `INSERT INTO post (user_id, content, allow_comment)
	          VALUES ($1, $2, $3)
	          RETURNING post_id, user_id, content, allow_comment, created_at`
	err := s.db.QueryRowxContext(ctx, query, input.UserID, input.Content, input.AllowComment).StructScan(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostgresStorage) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post

	query := `SELECT post_id, user_id, content, allow_comment, created_at FROM post WHERE post_id = $1`
	err := s.db.GetContext(ctx, &post, query, id)
	if err != nil {
		return nil, fmt.Errorf("post with id %d not found", id)
	}

	return &post, nil
}

func (s *PostgresStorage) GetPosts(ctx context.Context) ([]*models.Post, error) {
	var posts []*models.Post

	query := `SELECT post_id, user_id, content, allow_comment, created_at FROM post`
	err := s.db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, input model.NewComment) (*models.Comment, error) {
	if _, err := s.GetUserByID(ctx, input.UserID); err != nil {
		return nil, err
	}

	post, err := s.GetPostByID(ctx, input.PostID)
	if err != nil {
		return nil, err
	}

	if post.AllowComment != nil && !*post.AllowComment {
		return nil, fmt.Errorf("comments are disabled for post %d", input.PostID)
	}

	commentLevel := 1
	if input.ReplyCommentID != nil {
		var parentLevel int
		query := `SELECT comment_level FROM comment WHERE comment_id = $1`
		err := s.db.GetContext(ctx, &parentLevel, query, *input.ReplyCommentID)
		if err != nil {
			return nil, fmt.Errorf("comment with id %d not found", *input.ReplyCommentID)
		}
		commentLevel = parentLevel + 1
	}

	var comment models.Comment
	query := `INSERT INTO comment (reply_comment_id, comment_level, post_id, user_id, content)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING comment_id, reply_comment_id, comment_level, post_id, user_id, content, created_at`
	err = s.db.QueryRowxContext(ctx, query, input.ReplyCommentID, commentLevel, input.PostID, input.UserID, input.Content).StructScan(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *PostgresStorage) GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	var comments []*models.Comment

	query := `
		SELECT comment_id, post_id, reply_comment_id, comment_level, user_id, content, created_at
		FROM comment
		WHERE post_id = $1
		ORDER BY created_at ASC
	`
	err := s.db.SelectContext(ctx, &comments, query, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *PostgresStorage) GetCommentsByPostIDs(ctx context.Context, postIDs []int) (map[int][]*models.Comment, error) {
	result := make(map[int][]*models.Comment, len(postIDs))
	for _, postID := range postIDs {
		result[postID] = []*models.Comment{}
	}

	var comments []*models.Comment
	query := `
		SELECT comment_id, post_id, reply_comment_id, comment_level, user_id, content, created_at
		FROM comment
		WHERE post_id = ANY($1)
		ORDER BY created_at ASC
	`

	err := s.db.SelectContext(ctx, &comments, query, pq.Array(postIDs))
	if err != nil {
		return nil, err
	}

	for _, c := range comments {
		result[c.PostID] = append(result[c.PostID], c)
	}

	return result, nil
}
