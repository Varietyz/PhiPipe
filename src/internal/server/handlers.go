// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"regression-ci/pkg/types"
)

func (s *Server) healthCheck(c *gin.Context) {
	if err := s.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

func (s *Server) handleWebhook(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "webhook handling not implemented",
	})
}

func (s *Server) analyzeEndpoint(c *gin.Context) {
	var req types.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	result, err := s.detector.Analyze(req)
	if err != nil {
		log.Error().Err(err).Msg("analysis failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "analysis failed",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) getRepoConfig(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "config retrieval not implemented",
	})
}

func (s *Server) updateRepoConfig(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "config update not implemented",
	})
}