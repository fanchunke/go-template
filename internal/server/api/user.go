package api

import (
	"go-template/internal/server/cache"
	"go-template/internal/server/middleware"
	"go-template/internal/server/repository"
	"go-template/internal/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// UserAPI is the controller for user related requests.
type UserAPI struct {
	logger  *zap.Logger
	service service.UserService
}

// Get is the handler to get all users.
func (u *UserAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	logger := middleware.InjectedLogger(ctx, u.logger)
	logger.Info("start getting users")
	u.service.Get(ctx, "1")
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
	})
}

// NewUserAPI return an userAPI instance
func NewUserAPI(logger *zap.Logger, pool *redis.Pool, db *sqlx.DB) *UserAPI {
	repo := repository.NewUserRepo(db)
	cache := cache.NewUserCache(pool)
	service := service.NewUserService(repo, cache, logger)
	return &UserAPI{
		logger:  logger,
		service: service,
	}
}
