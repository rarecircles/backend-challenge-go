// Package api contains Rest APIs server and handlers.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/api/handler"
	"github.com/rarecircles/backend-challenge-go/internal/api/handler/tokengrp"
	"go.uber.org/zap"
)

// Config is configuration for server.
type Config struct {
	Log  *zap.Logger
	Addr string
}

// NewAPIServer creates http.Server that handle routes for the application.
func NewAPIServer(cfg *Config) *http.Server {
	r := gin.Default()

	r.Handle(http.MethodGet, "/healthcheck", handler.HealthCheck)

	th := tokengrp.NewHandler(cfg.Log)
	r.Handle(http.MethodGet, "/tokens", th.QueryTokens)

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	return &srv
}
