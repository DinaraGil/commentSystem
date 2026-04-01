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
	mu                      sync.RWMutex
	CommentPublishedChannel map[int]map[string]chan *models.Comment
}
