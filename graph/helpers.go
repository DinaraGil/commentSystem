package graph

import (
	"commentSystem/internal/models"
	"sort"
)

var subscriberCounter uint64

func getPaginatedComments(
	comments []*models.Comment,
	limit *int32,
	offset *int32,
	level *int32,
) []*models.Comment {
	children, roots := buildCommentTree(comments)
	sortComments(roots, children)
	selectedRoots := applyPagination(roots, limit, offset)

	return flattenTreeIterative(selectedRoots, children, level)
}

func buildCommentTree(comments []*models.Comment) (map[int][]*models.Comment, []*models.Comment) {
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

	return children, roots
}

func sortComments(roots []*models.Comment, children map[int][]*models.Comment) {
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].CreatedAt.After(roots[j].CreatedAt)
	})

	for parentID := range children {
		sort.Slice(children[parentID], func(i, j int) bool {
			return children[parentID][i].CreatedAt.After(children[parentID][j].CreatedAt)
		})
	}
}

func applyPagination(
	roots []*models.Comment,
	limit *int32,
	offset *int32,
) []*models.Comment {
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

	return roots[o:end]
}

func flattenTreeIterative(
	roots []*models.Comment,
	children map[int][]*models.Comment,
	level *int32,
) []*models.Comment {
	result := make([]*models.Comment, 0)

	stack := make([]*models.Comment, 0, len(roots))

	for i := len(roots) - 1; i >= 0; i-- {
		stack = append(stack, roots[i])
	}

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		if level != nil && int32(current.CommentLevel) > *level {
			continue
		}

		result = append(result, current)

		currentChildren := children[current.ID]
		for i := len(currentChildren) - 1; i >= 0; i-- {
			stack = append(stack, currentChildren[i])
		}
	}

	return result
}
