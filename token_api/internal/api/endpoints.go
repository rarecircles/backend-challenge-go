package api

// Token API endpoints' definitions. It's here where routes are defined and HTTP requests
// from third parties are processed.

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/dataloader"
	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/types"
)

type tokenAPI struct {
	tokenData *dataloader.ProcessedTokens
}

func newtokenAPI(data *dataloader.ProcessedTokens) (tokenAPI, error) {
	return tokenAPI{tokenData: data}, nil
}

func (m *tokenAPI) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	})
}
func (m tokenAPI) getTokens(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		m.writeJSONError(w, 400, "q is required")
		return
	}
	indexices := m.tokenData.TokenNameSuffixArray.Lookup([]byte(q), -1)
	matches := []types.Token{}
	for _, index := range indexices {
		tokenName := m.tokenData.GetStringFromIndex(index)
		token, ok := m.tokenData.TokenMap[tokenName]
		if !ok {
			log.Errorf("failed to get %s from token map", tokenName)
		}
		if !slices.Contains(matches, token) {
			matches = append(matches, token)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	response := types.TokensResult{Tokens: matches}
	json.NewEncoder(w).Encode(response)
}

// TokenMux creates a HTTP handler for talking to the Metrics API.
func TokenMux(data *dataloader.ProcessedTokens) (*mux.Router, error) {
	taMux := mux.NewRouter()

	endpoints, err := newtokenAPI(data)
	if err != nil {
		return nil, fmt.Errorf("newtokenAPI(): %w", err)
	}

	// Endpoints for the Metrics API
	taMux.HandleFunc("/tokens", endpoints.getTokens).Methods("GET")
	return taMux, nil
}
