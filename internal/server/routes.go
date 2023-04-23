package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.engine

	router.GET("/", s.handlers.Up)
	router.GET("/tokens", s.handlers.SearchTokens)

	return router
}
