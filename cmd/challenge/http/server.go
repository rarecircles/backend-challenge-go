package http

import (
	"net/http"

	"github.com/go-chi/chi"
	endpoint "github.com/rarecircles/backend-challenge-go/cmd/challenge/http/endpoints"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	server *http.Server
}

func (s *Server) setAllRoutes() error {
	s.router.Get("/tokens", endpoint.Tokens)

	return nil
}

func CreateServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
		server: &http.Server{Addr: ":8080"},
	}

	err := s.setAllRoutes()
	if err != nil {
		zlog.Error("Error setting routes", zap.String("error", err.Error()))
	}

	s.server.Handler = s.router
	return s

}
func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		zlog.Error("Error starting http server", zap.String("error", err.Error()))
	}
}
