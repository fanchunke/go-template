package api

import (
	"go-template/internal/errno"
	"go-template/internal/log"
	"go-template/internal/server/cache"
	"go-template/internal/server/repository"
	"go-template/internal/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// UserAPI is the controller for user related requests.
type UserAPI struct {
	service service.UserService
}

// Get is the handler to get all users.
func (u *UserAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)
	logger.Info("start getting users")
	if _, err := u.service.Get(ctx, "1"); err != nil {
		c.JSON(http.StatusInternalServerError, errno.ErrServer.WithError(err))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
	})
}

// NewUserAPI return an userAPI instance
func NewUserAPI(pool *redis.Pool, db *sqlx.DB) *UserAPI {
	repo := repository.NewUserRepo(db)
	cache := cache.NewUserCache(pool)
	service := service.NewUserService(repo, cache)
	return &UserAPI{
		service: service,
	}
}
