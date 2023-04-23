package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/handlers"
	"go.uber.org/zap"
)

type Server struct {
	r *gin.Engine
	h handlers.Interface
	l *zap.Logger
}

func New(l *zap.Logger, h handlers.Interface) *Server {
	router := gin.Default()
	return &Server{
		r: router,
		h: h,
		l: l,
	}
}

func (s *Server) Run(baseUrl, port string) error {
	r := s.Routes()

	return r.Run(baseUrl + port)
}
