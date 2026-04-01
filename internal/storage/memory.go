package storage

import (
	"commentSystem/graph/model"
	"commentSystem/internal/models"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
	"unicode/utf8"
)

type MemoryStorage struct {
	mu sync.RWMutex

	users    map[int]*models.User
	posts    map[int]*models.Post
	comments map[int]*models.Comment

	nextUserID    int
	nextPostID    int
	nextCommentID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:         make(map[int]*models.User),
		posts:         make(map[int]*models.Post),
		comments:      make(map[int]*models.Comment),
		nextUserID:    1,
		nextPostID:    1,
		nextCommentID: 1,
	}
}

func (s *MemoryStorage) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &models.User{
		ID:       s.nextUserID,
		UserName: input.Username,
	}
	s.users[user.ID] = user
	s.nextUserID++

	return user, nil
}

func (s *MemoryStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (s *MemoryStorage) CreatePost(ctx context.Context, input model.NewPost) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[input.UserID]; !ok {
		return nil, fmt.Errorf("user with id %d not found", input.UserID)
	}

	allowComment := input.AllowComment
	post := &models.Post{
		ID:           s.nextPostID,
		UserID:       input.UserID,
		Content:      input.Content,
		AllowComment: allowComment,
		CreatedAt:    time.Now(),
	}

	s.posts[post.ID] = post
	s.nextPostID++

	return post, nil
}

func (s *MemoryStorage) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", id)
	}
	return post, nil
}

func (s *MemoryStorage) GetPosts(ctx context.Context) ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	return posts, nil
}

func (s *MemoryStorage) CreateComment(ctx context.Context, input model.NewComment) (*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if utf8.RuneCountInString(input.Content) > 2000 {
		return nil, fmt.Errorf("the comment should contain less than 2000")
	}
	if _, ok := s.users[input.UserID]; !ok {
		return nil, fmt.Errorf("user with id %d not found", input.UserID)
	}

	post, ok := s.posts[input.PostID]
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", input.PostID)
	}

	if post.AllowComment != nil && !*post.AllowComment {
		return nil, fmt.Errorf("comments are disabled for post %d", input.PostID)
	}

	commentLevel := 1
	if input.ReplyCommentID != nil {
		parent, ok := s.comments[*input.ReplyCommentID]
		if !ok {
			return nil, fmt.Errorf("comment with id %d not found", *input.ReplyCommentID)
		}
		commentLevel = parent.CommentLevel + 1
	}

	comment := &models.Comment{
		ID:             s.nextCommentID,
		ReplyCommentID: input.ReplyCommentID,
		CommentLevel:   commentLevel,
		PostID:         input.PostID,
		UserID:         input.UserID,
		Content:        input.Content,
		CreatedAt:      time.Now(),
		Replies:        nil,
	}

	s.comments[comment.ID] = comment
	s.nextCommentID++

	return comment, nil
}

func (s *MemoryStorage) GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments := make([]*models.Comment, 0)
	for _, comment := range s.comments {
		if comment.PostID == postID {
			c := *comment
			c.Replies = nil
			comments = append(comments, &c)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	return comments, nil
}
