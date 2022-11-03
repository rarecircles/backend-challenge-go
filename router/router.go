package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rarecircles/backend-challenge-go/handler"
	"github.com/rarecircles/backend-challenge-go/service"
)

type Options struct {
	Dependencies *service.ServerDependencies
}

func InitRouter(opt Options) http.Handler {
	h := http.NewServeMux()
	router := mux.NewRouter()

	router.HandleFunc("/ping", handler.Ping()).Methods("GET")

	router.HandleFunc("/tokens", handler.GetTokensInfo(opt.Dependencies.TokenService)).Methods("GET")

	h.Handle("/", router)
	return h
}
