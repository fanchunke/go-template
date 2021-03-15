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

	"go-template/internal/server/config"
	"go-template/internal/server/router"
	"go-template/internal/version"
)

// Config is app configuration
type Config config.Config

// Server is a HTTP server
type Server struct {
	logger *zap.Logger
	config *Config
	router http.Handler
}

// NewServer return a HTTP server
func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	srv := &Server{
		config: config,
		logger: logger,
	}
	return srv, nil
}

// Run starts HTTP server and watches channel to determine whether to stop server gracefully.
func (s *Server) Run(stopCh <-chan struct{}) {
	go s.startMetricsServer()

	// // init app context, include create redis pool, connect database, etc.
	ticker := time.NewTicker(30 * time.Second)
	pool := s.startCachePool(ticker, stopCh)
	db, err := s.connectDatabase()
	if err != nil {
		panic(err)
	}

	// register http handlers
	s.registerHandlers(pool, db)

	// create http server
	srv := s.startServer()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HTTPServerShutdownTimeout)
	defer cancel()

	// close cache pool
	if pool != nil {
		_ = pool.Close()
	}

	s.logger.Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.config.HTTPServerShutdownTimeout))

	// determine if the http server was started
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}
}

func (s *Server) registerHandlers(pool *redis.Pool, db *sqlx.DB) {
	s.router = router.New(s.logger, pool, db)
}

func (s *Server) startServer() *http.Server {
	// determine if the port is specified
	c := s.config
	if c.Port == "0" {
		return nil
	}

	srv := &http.Server{
		Addr:         ":" + c.Port,
		WriteTimeout: c.HTTPServerTimeout,
		ReadTimeout:  c.HTTPServerTimeout,
		IdleTimeout:  2 * c.HTTPServerShutdownTimeout,
		Handler:      s.router,
	}

	// start the server in the background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	return srv
}

func (s *Server) startMetricsServer() {
	if s.config.PortMetrics > 0 {
		mux := http.DefaultServeMux
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%v", s.config.PortMetrics),
			Handler: mux,
		}
		srv.ListenAndServe()
	}
}

func (s *Server) startCachePool(ticker *time.Ticker, stopCh <-chan struct{}) *redis.Pool {
	c := s.config.Redis
	if c == nil {
		return nil
	}

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

	setVersion := func() {
		conn := pool.Get()
		defer conn.Close()
		if _, err := conn.Do("SET", s.config.HostName, version.VERSION, "EX", 60); err != nil {
			s.logger.Warn("cache server is offline", zap.Error(err), zap.Any("server", c))
		}
	}

	// set version on a schedule
	go func() {
		setVersion()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				setVersion()
			}
		}
	}()

	return pool
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
