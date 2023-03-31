package handler

import (
	"encoding/json"
	"github.com/rarecircles/backend-challenge-go/cmd/dao"
	"github.com/rarecircles/backend-challenge-go/cmd/service"
	"net/http"
)

func GetTokens(DAO dao.DaoInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")

		tokensDTO, err := service.GetTokens(DAO, q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorMsg := InitializeError(err.Error())
			w.Write(errorMsg)
			return
		}

		jsonResp, err := json.Marshal(tokensDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorMsg := InitializeError(err.Error())
			w.Write(errorMsg)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return
	}
}
