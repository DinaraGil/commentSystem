package storage

import (
	"commentSystem/graph/model"
	"context"
	"testing"
)

func TestMemoryStorage_CreateUser(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, err := store.CreateUser(ctx, model.NewUser{
		Username: "alice",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected user ID = 1, got %d", user.ID)
	}

	if user.UserName != "alice" {
		t.Fatalf("expected username = alice, got %s", user.UserName)
	}
}

func TestMemoryStorage_CreatePost_Success(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, err := store.CreateUser(ctx, model.NewUser{
		Username: "alice",
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	allow := true
	post, err := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "hello post",
		AllowComment: &allow,
	})
	if err != nil {
		t.Fatalf("CreatePost returned error: %v", err)
	}

	if post.ID != 1 {
		t.Fatalf("expected post ID = 1, got %d", post.ID)
	}

	if post.UserID != user.ID {
		t.Fatalf("expected post.UserId = %d, got %d", user.ID, post.UserID)
	}

	if post.Content != "hello post" {
		t.Fatalf("expected content = hello post, got %s", post.Content)
	}
}

func TestMemoryStorage_CreatePost_UserNotFound(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	allow := true
	_, err := store.CreatePost(ctx, model.NewPost{
		UserID:       999,
		Content:      "hello post",
		AllowComment: &allow,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMemoryStorage_CreateComment_Success(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	comment, err := store.CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "first comment",
	})
	if err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}

	if comment.ID != 1 {
		t.Fatalf("expected comment ID = 1, got %d", comment.ID)
	}

	if comment.CommentLevel != 1 {
		t.Fatalf("expected level = 1, got %d", comment.CommentLevel)
	}

	if comment.PostID != post.ID {
		t.Fatalf("expected postID = %d, got %d", post.ID, comment.PostID)
	}
}

func TestMemoryStorage_CreateComment_Disabled(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := false
	post, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	_, err := store.CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "comment",
	})
	if err == nil {
		t.Fatal("expected error when comments are disabled, got nil")
	}
}

func TestMemoryStorage_CreateComment_ReplyLevel(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	parent, err := store.CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "parent",
	})
	if err != nil {
		t.Fatalf("CreateComment parent returned error: %v", err)
	}

	reply, err := store.CreateComment(ctx, model.NewComment{
		ReplyCommentID: &parent.ID,
		PostID:         post.ID,
		UserID:         user.ID,
		Content:        "reply",
	})
	if err != nil {
		t.Fatalf("CreateComment reply returned error: %v", err)
	}

	if reply.CommentLevel != 2 {
		t.Fatalf("expected reply level = 2, got %d", reply.CommentLevel)
	}
}

func TestMemoryStorage_CreateComment_ParentNotFound(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	parentID := 999
	_, err := store.CreateComment(ctx, model.NewComment{
		ReplyCommentID: &parentID,
		PostID:         post.ID,
		UserID:         user.ID,
		Content:        "reply",
	})
	if err == nil {
		t.Fatal("expected error for missing parent comment, got nil")
	}
}

func TestMemoryStorage_GetCommentsByPostID(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post1, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post1",
		AllowComment: &allow,
	})
	post2, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post2",
		AllowComment: &allow,
	})

	_, _ = store.CreateComment(ctx, model.NewComment{
		PostID:  post1.ID,
		UserID:  user.ID,
		Content: "comment 1",
	})
	_, _ = store.CreateComment(ctx, model.NewComment{
		PostID:  post1.ID,
		UserID:  user.ID,
		Content: "comment 2",
	})
	_, _ = store.CreateComment(ctx, model.NewComment{
		PostID:  post2.ID,
		UserID:  user.ID,
		Content: "comment 3",
	})

	comments, err := store.GetCommentsByPostID(ctx, post1.ID)
	if err != nil {
		t.Fatalf("GetCommentsByPostID returned error: %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}

	for _, c := range comments {
		if c.PostID != post1.ID {
			t.Fatalf("expected PostID = %d, got %d", post1.ID, c.PostID)
		}
	}
}

func TestMemoryStorage_CreateComment_ContentTooLong(t *testing.T) {
	store := NewMemoryStorage()
	ctx := context.Background()

	user, _ := store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post, _ := store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	longText := make([]byte, 2001)
	for i := range longText {
		longText[i] = 'a'
	}

	_, err := store.CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: string(longText),
	})
	if err == nil {
		t.Fatal("expected error for too long comment, got nil")
	}
}
