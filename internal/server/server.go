package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/handlers"
	"go.uber.org/zap"
)

type Server struct {
	engine   *gin.Engine
	handlers handlers.I
	logger   *zap.Logger
}

func New(l *zap.Logger, h handlers.I) *Server {
	router := gin.Default()
	return &Server{
		engine:   router,
		handlers: h,
		logger:   l,
	}
}

func (s *Server) Run(baseUrl, port string) error {
	r := s.Routes()

	return r.Run(baseUrl + port)
}
