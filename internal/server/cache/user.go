package cache

import (
	"go-template/internal/server/model"
	"log"

	"github.com/gomodule/redigo/redis"
)

type userCache struct {
	pool *redis.Pool
}

// NewUserCache creates an UserCache instance.
func NewUserCache(pool *redis.Pool) UserCache {
	return &userCache{
		pool: pool,
	}
}

func (c *userCache) Get(userID string) (*model.UserCache, error) {
	conn := c.pool.Get()
	defer conn.Close()

	var (
		user model.UserCache
	)
	values, err := redis.Values(conn.Do("HGETALL", userID))
	if err != nil {
		return nil, err
	}
	if err = redis.ScanStruct(values, &user); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func (c *userCache) Set(userID string, user *model.UserCache) error {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("HMSET", redis.Args{}.Add(userID).AddFlat(user)...); err != nil {
		return err
	}
	return nil
}
