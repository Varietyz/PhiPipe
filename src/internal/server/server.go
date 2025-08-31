// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"regression-ci/internal/config"
	"regression-ci/internal/regression"
	"regression-ci/internal/github"
)

type Server struct {
	db       *sqlx.DB
	config   *config.Config
	detector *regression.Detector
	github   *github.Client
	router   *gin.Engine
	server   *http.Server
}

func New(db *sqlx.DB, cfg *config.Config) *Server {
	s := &Server{
		db:       db,
		config:   cfg,
		detector: regression.New(db, cfg.Detection),
		github:   github.New(cfg.GitHub),
	}

	s.setupRoutes()
	s.server = &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      s.router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) setupRoutes() {
	s.router = gin.New()
	s.router.Use(gin.Recovery())
	s.router.Use(s.loggingMiddleware())

	s.router.GET("/health", s.healthCheck)
	s.router.POST("/webhook", s.handleWebhook)
	s.router.POST("/analyze", s.analyzeEndpoint)
	s.router.GET("/config/:repo", s.getRepoConfig)
	s.router.PUT("/config/:repo", s.updateRepoConfig)
}

func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info().
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Msg("request completed")
	}
}