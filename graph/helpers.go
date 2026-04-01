package graph

import "commentSystem/internal/models"

//
//import (
//	"commentSystem/internal/models"
//	"context"
//	"database/sql"
//	"fmt"
//)
//
//func (r *Resolver) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
//	var user models.User
//	query := `SELECT person_id, username FROM person WHERE person_id=$1`
//	err := r.DB.Get(&user, query, userID)
//	if err == sql.ErrNoRows {
//		return nil, fmt.Errorf("user with id %d not found", userID)
//	}
//	return &user, err
//}
//
//func (r *Resolver) GetPostByID(ctx context.Context, postId int) (*models.Post, error) {
//	var post models.Post
//
//	query := `SELECT * FROM post WHERE post_id=$1`
//	err := r.DB.Get(&post, query, postId)
//	if err == sql.ErrNoRows {
//		return nil, fmt.Errorf("post with id %d not found", postId)
//	}
//	return &post, err
//}
//
//func buildCommentTree(comments []*models.Comment) []*models.Comment {
//	commentMap := make(map[int]*models.Comment)
//	var roots []*models.Comment
//
//	for _, c := range comments {
//		commentMap[c.ID] = c
//	}
//
//	for _, c := range comments {
//		if c.ReplyCommentID != nil {
//			parent := commentMap[*c.ReplyCommentID]
//			if parent != nil {
//				parent.Replies = append(parent.Replies, c)
//			}
//		} else {
//			roots = append(roots, c)
//		}
//	}
//
//	return roots
//}

import "sort"

var subscriberCounter uint64

func flattenCommentsAsTree(
	comments []*models.Comment,
	limit *int32,
	offset *int32,
	level *int32,
) []*models.Comment {
	children := make(map[int][]*models.Comment)
	roots := make([]*models.Comment, 0)

	for _, c := range comments {
		if c.ReplyCommentID == nil {
			roots = append(roots, c)
			continue
		}

		parentID := *c.ReplyCommentID
		children[parentID] = append(children[parentID], c)
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].CreatedAt.After(roots[j].CreatedAt)
	})

	for parentID := range children {
		sort.Slice(children[parentID], func(i, j int) bool {
			return children[parentID][i].CreatedAt.After(children[parentID][j].CreatedAt)
		})
	}

	var l int32 = int32(len(roots))
	var o int32 = 0

	if limit != nil && *limit >= 0 {
		l = *limit
	}
	if offset != nil && *offset >= 0 {
		o = *offset
	}

	if o >= int32(len(roots)) {
		return []*models.Comment{}
	}

	end := o + l
	if end > int32(len(roots)) {
		end = int32(len(roots))
	}

	selectedRoots := roots[o:end]

	result := make([]*models.Comment, 0)

	var walk func(c *models.Comment)
	walk = func(c *models.Comment) {
		if level != nil && int32(c.CommentLevel) > *level {
			return
		}

		result = append(result, c)

		for _, child := range children[c.ID] {
			walk(child)
		}
	}

	for _, root := range selectedRoots {
		walk(root)
	}

	return result
}
