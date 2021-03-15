package service

import (
	"context"
	"go-template/internal/server/model"
	"go-template/internal/server/repository"
)

type bookService struct {
	repo repository.BookRepo
}

func (s *bookService) Get(ctx context.Context, bookID string) (*model.Book, error) {
	return s.repo.Get(bookID)
}

// NewBookService returns a BookService instance.
func NewBookService(repo repository.BookRepo) BookService {
	return &bookService{
		repo: repo,
	}
}
