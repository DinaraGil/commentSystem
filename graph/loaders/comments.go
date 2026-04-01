package loaders

import (
	"commentSystem/internal/models"
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type commentReader struct {
	db *sql.DB
}

func (r *commentReader) getCommentsByPostIDs(ctx context.Context, postIDs []int) ([][]*models.Comment, []error) {
	query := `
		SELECT comment_id, post_id, reply_comment_id, comment_level, user_id, content, created_at
		FROM comment
		WHERE post_id = ANY($1)
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(postIDs))
	if err != nil {
		return nil, []error{err}
	}
	defer rows.Close()

	commentMap := make(map[int][]*models.Comment, len(postIDs))

	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ReplyCommentID,
			&c.CommentLevel,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, []error{err}
		}

		comment := c
		commentMap[c.PostID] = append(commentMap[c.PostID], &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	results := make([][]*models.Comment, len(postIDs))
	errs := make([]error, len(postIDs))

	for i, postID := range postIDs {
		if comments, ok := commentMap[postID]; ok {
			results[i] = comments
		} else {
			results[i] = []*models.Comment{}
		}
		errs[i] = nil
	}

	return results, errs
}
