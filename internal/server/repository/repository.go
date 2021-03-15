package repository

import "go-template/internal/server/model"

// UserRepo is an interface to access user table
type UserRepo interface {
	Get(userID string) (*model.User, error)
}

// BookRepo is an interface to access book table
type BookRepo interface {
	Get(bookID string) (*model.Book, error)
}
