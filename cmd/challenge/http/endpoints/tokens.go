package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/rarecircles/backend-challenge-go/cmd/challenge/sql"
	"github.com/rarecircles/backend-challenge-go/cmd/challenge/types"
	"go.uber.org/zap"
)

func Tokens(w http.ResponseWriter, r *http.Request) {
	query := string(r.URL.Query()["q"][0])

	tokens, err := sql.TokensFetch(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	rsp, err := json.MarshalIndent(types.TokenQueryResponse{tokens}, "", "  ")
	if err != nil {
		zlog.Error("Error fetching tokens by title", zap.String("title", query), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
}
