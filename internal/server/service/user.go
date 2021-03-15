package service

import (
	"context"
	"go-template/internal/server/cache"
	"go-template/internal/server/middleware"
	"go-template/internal/server/model"
	"go-template/internal/server/repository"

	"go.uber.org/zap"
)

type userService struct {
	repo   repository.UserRepo
	cache  cache.UserCache
	logger *zap.Logger
}

// NewUserService returns an UserService instance.
func NewUserService(repo repository.UserRepo, cache cache.UserCache, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *userService) Get(ctx context.Context, userID string) (*model.User, error) {
	logger := middleware.InjectedLogger(ctx, s.logger)
	logger.Info("start userService.Get")
	return s.repo.Get(userID)
}
