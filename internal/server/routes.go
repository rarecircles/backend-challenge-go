package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.r

	router.GET("/", s.h.Up)
	router.GET("/tokens", s.h.SearchTokens)

	return router
}
