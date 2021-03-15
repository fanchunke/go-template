package service

import (
	"context"
	"go-template/internal/server/model"
)

// UserService interface
type UserService interface {
	Get(ctx context.Context, userID string) (*model.User, error)
}

// BookService interface
type BookService interface {
	Get(ctx context.Context, bookID string) (*model.Book, error)
}
