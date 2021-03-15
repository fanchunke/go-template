package cache

import "go-template/internal/server/model"

// UserCache is an interface to get user info from cache.
type UserCache interface {
	Get(userID string) (*model.UserCache, error)
	Set(userID string, user *model.UserCache) error
}
