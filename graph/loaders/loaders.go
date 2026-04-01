package loaders

import (
	"commentSystem/internal/models"
	"context"
	"database/sql"
	"time"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const loadersKey ctxKey = "dataloaders"

type Loaders struct {
	CommentsLoader *dataloadgen.Loader[int, []*models.Comment]
}

func NewLoaders(db *sql.DB) *Loaders {
	cr := &commentReader{db: db}

	return &Loaders{
		CommentsLoader: dataloadgen.NewLoader(
			cr.getCommentsByPostIDs,
			dataloadgen.WithWait(1*time.Millisecond),
		),
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetCommentsByPostID загружает комментарии одного поста через batch loader
func GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	loaders := For(ctx)
	return loaders.CommentsLoader.Load(ctx, postID)
}

// GetCommentsByPostIDs — опционально, если где-то понадобится массовая загрузка
func GetCommentsByPostIDs(ctx context.Context, postIDs []int) ([][]*models.Comment, error) {
	loaders := For(ctx)
	return loaders.CommentsLoader.LoadAll(ctx, postIDs)
}
