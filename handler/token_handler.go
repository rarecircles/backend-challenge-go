package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rarecircles/backend-challenge-go/service"
	"github.com/rarecircles/backend-challenge-go/views"
)

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal("pong")
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

func GetTokensInfo(service service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("q")
		tokenTitle := strings.ToLower(param)
		if len(tokenTitle) == 0 {
			fmt.Println("need query param, q")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := service.GetTokensInfo(tokenTitle)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := views.Resp{Tokens: data}
		body, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}
