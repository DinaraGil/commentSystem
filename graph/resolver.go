package graph

import (
	"commentSystem/internal/models"
	"sync"

	"github.com/jmoiron/sqlx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	DB                      *sqlx.DB
	mu                      sync.Mutex
	CommentPublishedChannel map[int][]chan *models.Comment
}

//var commentPublishedChannel map[int][]chan *models.Comment
//
//func init() {
//	commentPublishedChannel = map[int][]chan *models.Comment{}
//}
