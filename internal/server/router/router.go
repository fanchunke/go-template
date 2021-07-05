package router

import (
	"go-template/internal/server/api"
	"go-template/internal/server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// New returns a http.Handler
func New(pool *redis.Pool, db *sqlx.DB) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// RequestID middleware must be registered at the beginning.
	r.Use(middleware.RequestID())
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Prometheus())
	r.Use(middleware.Version())

	userAPI := api.NewUserAPI(pool, db)
	// r.GET("/", api.Index.Healthy(env))
	r.GET("/users", userAPI.Get)
	return r
}
