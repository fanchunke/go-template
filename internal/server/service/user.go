package service

import (
	"context"
	"go-template/internal/log"
	"go-template/internal/server/cache"
	"go-template/internal/server/model"
	"go-template/internal/server/repository"
)

type userService struct {
	repo  repository.UserRepo
	cache cache.UserCache
}

// NewUserService returns an UserService instance.
func NewUserService(repo repository.UserRepo, cache cache.UserCache) UserService {
	return &userService{
		repo:  repo,
		cache: cache,
	}
}

func (s *userService) Get(ctx context.Context, userID string) (*model.User, error) {
	logger := log.Ctx(ctx)
	logger.Info("start userService.Get")
	// return s.repo.Get(userID)
	userCache, err := s.cache.Get(userID)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:   userCache.ID,
		Name: userCache.Name,
	}
	return user, nil
}
