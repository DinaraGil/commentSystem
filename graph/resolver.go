package graph

import (
	"commentSystem/internal/models"
	"commentSystem/internal/storage"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	//DB                      *sqlx.DB
	Store                   storage.Storage
	mu                      sync.Mutex
	CommentPublishedChannel map[int][]chan *models.Comment
}

//var commentPublishedChannel map[int][]chan *models.Comment
//
//func init() {
//	commentPublishedChannel = map[int][]chan *models.Comment{}
//}
