package graph

import (
	"commentSystem/graph/model"
	"commentSystem/internal/models"
	"commentSystem/internal/storage"
	"context"
	"sync"
	"testing"
	"time"
)

func newTestResolver() *Resolver {
	return &Resolver{
		Store:                   storage.NewMemoryStorage(),
		CommentPublishedChannel: make(map[int]map[string]chan *models.Comment),
		mu:                      sync.RWMutex{},
	}
}

func TestSubscription_MultipleSubscribersReceiveComment(t *testing.T) {
	r := newTestResolver()
	ctx := context.Background()

	user, err := r.Store.CreateUser(ctx, model.NewUser{Username: "alice"})
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}

	allow := true
	post, err := r.Store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})
	if err != nil {
		t.Fatalf("CreatePost error: %v", err)
	}

	subCtx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	subCtx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	ch1, err := (&subscriptionResolver{r}).CommentPublished(subCtx1, post.ID)
	if err != nil {
		t.Fatalf("CommentPublished #1 error: %v", err)
	}

	ch2, err := (&subscriptionResolver{r}).CommentPublished(subCtx2, post.ID)
	if err != nil {
		t.Fatalf("CommentPublished #2 error: %v", err)
	}

	created, err := (&mutationResolver{r}).CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "new comment",
	})
	if err != nil {
		t.Fatalf("CreateComment error: %v", err)
	}

	select {
	case got := <-ch1:
		if got.ID != created.ID {
			t.Fatalf("subscriber #1 got wrong comment ID: want %d, got %d", created.ID, got.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber #1 did not receive comment")
	}

	select {
	case got := <-ch2:
		if got.ID != created.ID {
			t.Fatalf("subscriber #2 got wrong comment ID: want %d, got %d", created.ID, got.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber #2 did not receive comment")
	}
}

func TestSubscription_OtherPostDoesNotReceiveComment(t *testing.T) {
	r := newTestResolver()
	ctx := context.Background()

	user, err := r.Store.CreateUser(ctx, model.NewUser{Username: "alice"})
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}

	allow := true
	post1, _ := r.Store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post1",
		AllowComment: &allow,
	})
	post2, _ := r.Store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post2",
		AllowComment: &allow,
	})

	subCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := (&subscriptionResolver{r}).CommentPublished(subCtx, post2.ID)
	if err != nil {
		t.Fatalf("CommentPublished error: %v", err)
	}

	_, err = (&mutationResolver{r}).CreateComment(ctx, model.NewComment{
		PostID:  post1.ID,
		UserID:  user.ID,
		Content: "comment for post1",
	})
	if err != nil {
		t.Fatalf("CreateComment error: %v", err)
	}

	select {
	case got := <-ch:
		t.Fatalf("unexpected comment received for another post: %+v", got)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestSubscription_UnsubscribeOnContextDone(t *testing.T) {
	r := newTestResolver()
	ctx := context.Background()

	user, err := r.Store.CreateUser(ctx, model.NewUser{Username: "alice"})
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}

	allow := true
	post, err := r.Store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})
	if err != nil {
		t.Fatalf("CreatePost error: %v", err)
	}

	subCtx, cancel := context.WithCancel(context.Background())

	_, err = (&subscriptionResolver{r}).CommentPublished(subCtx, post.ID)
	if err != nil {
		t.Fatalf("CommentPublished error: %v", err)
	}

	cancel()
	time.Sleep(50 * time.Millisecond)

	r.mu.RLock()
	defer r.mu.RUnlock()

	subs := r.CommentPublishedChannel[post.ID]
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscribers after cancel, got %d", len(subs))
	}
}

func TestPostResolver_Comments_ReturnsTree(t *testing.T) {
	r := newTestResolver()
	ctx := context.Background()

	user, _ := r.Store.CreateUser(ctx, model.NewUser{Username: "alice"})
	allow := true
	post, _ := r.Store.CreatePost(ctx, model.NewPost{
		UserID:       user.ID,
		Content:      "post",
		AllowComment: &allow,
	})

	parent, _ := r.Store.CreateComment(ctx, model.NewComment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "parent",
	})

	_, _ = r.Store.CreateComment(ctx, model.NewComment{
		ReplyCommentID: &parent.ID,
		PostID:         post.ID,
		UserID:         user.ID,
		Content:        "child",
	})

	limit := int32(10)
	offset := int32(0)
	level := int32(10)

	comments, err := (&postResolver{r}).Comments(ctx, post, &limit, &offset, &level)
	if err != nil {
		t.Fatalf("Comments error: %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}
