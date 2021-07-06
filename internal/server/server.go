package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"go-template/internal/config"
	"go-template/internal/server/router"
)

// Server is a HTTP server
type Server struct {
	config *config.Config
	router http.Handler
}

// NewServer return a HTTP server
func NewServer(config *config.Config) (*Server, error) {
	srv := &Server{
		config: config,
	}
	return srv, nil
}

// Run starts HTTP server and watches channel to determine whether to stop server gracefully.
func (s *Server) Run(stopCh <-chan struct{}) {
	go s.startMetricsServer()

	// create redis pool, connect database, etc.
	pool, err := s.startCachePool()
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("connect to redis failed: %v\n", err.Error()))
	}

	db, err := s.connectDatabase()
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("connect to database failed: %v\n", err.Error()))
	}

	// register http handlers
	s.registerHandlers(pool, db)

	// create http server
	srv := s.startServer()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HTTP.HTTPServerShutdownTimeout)
	defer cancel()

	// close cache pool
	if pool != nil {
		_ = pool.Close()
	}

	zap.L().Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.config.HTTP.HTTPServerShutdownTimeout))

	// determine if the http server was started
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			zap.L().Warn("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}
}

func (s *Server) registerHandlers(pool *redis.Pool, db *sqlx.DB) {
	s.router = router.New(pool, db)
}

func (s *Server) startServer() *http.Server {
	// determine if the port is specified
	c := s.config
	if c.HTTP.Port == "0" {
		return nil
	}

	srv := &http.Server{
		Addr:         ":" + c.HTTP.Port,
		WriteTimeout: c.HTTP.HTTPServerTimeout,
		ReadTimeout:  c.HTTP.HTTPServerTimeout,
		IdleTimeout:  2 * c.HTTP.HTTPServerShutdownTimeout,
		Handler:      s.router,
	}

	// start the server in the background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			zap.L().Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	return srv
}

func (s *Server) startMetricsServer() {
	if s.config.HTTP.PortMetrics > 0 {
		mux := http.DefaultServeMux
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%v", s.config.HTTP.PortMetrics),
			Handler: mux,
		}
		srv.ListenAndServe()
	}
}

func (s *Server) startCachePool() (*redis.Pool, error) {
	c := s.config.Redis

	pool := &redis.Pool{
		MaxIdle:     c.MaxIdle,
		IdleTimeout: c.IdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp", fmt.Sprintf("%s:%s", c.Host, c.Port),
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.DB),
				redis.DialConnectTimeout(c.ConnectTimeout),
				redis.DialReadTimeout(c.ReadTimeout),
				redis.DialWriteTimeout(c.WriteTimeout),
			)
		},
	}

	// test connection pool
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return pool, err
	}

	return pool, nil
}

func (s *Server) connectDatabase() (*sqlx.DB, error) {
	c := s.config.Database
	config := &mysql.Config{
		User:                 c.User,
		Passwd:               c.Password,
		DBName:               c.DBName,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", c.Host, c.Port),
		Loc:                  time.Local,
		Timeout:              c.Timeout,
		ReadTimeout:          c.ReadTimeout,
		WriteTimeout:         c.WriteTimeout,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}
	return sqlx.Connect("mysql", config.FormatDSN())
}
