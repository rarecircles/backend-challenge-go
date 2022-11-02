// Package search
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"go.uber.org/zap"
)

const (
	TokensIndex = "tokens"
)

type SearchService interface {
	SearchToken(ctx context.Context, name string) ([]eth.Token, error)
}

type searchService struct {
	log      *zap.Logger
	esClient *elasticsearch.Client
}

func NewSearchService(log *zap.Logger, esClient *elasticsearch.Client) SearchService {
	return &searchService{
		log:      log,
		esClient: esClient,
	}
}

type SearchTokenResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []*struct {
			Source *eth.Token `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *searchService) SearchToken(ctx context.Context, name string) ([]eth.Token, error) {
	var searchBuffer bytes.Buffer
	search := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"name": map[string]any{
					"query":     name,
					"fuzziness": "AUTO",
				},
			},
		},
	}

	if err := json.NewEncoder(&searchBuffer).Encode(search); err != nil {
		return nil, fmt.Errorf("encoding json: %w", err)
	}

	resp, err := s.esClient.Search(
		s.esClient.Search.WithContext(ctx),
		s.esClient.Search.WithIndex(TokensIndex),
		s.esClient.Search.WithBody(&searchBuffer),
		s.esClient.Search.WithTrackTotalHits(true),
		s.esClient.Search.WithPretty(),
	)

	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	var searchReponse = SearchTokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&searchReponse); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}

	ethTokens := []eth.Token{}
	if searchReponse.Hits.Total.Value > 0 {
		ethTokens = make([]eth.Token, len(searchReponse.Hits.Hits))
		for i, hit := range searchReponse.Hits.Hits {
			ethTokens[i].Name = hit.Source.Name
			ethTokens[i].Decimals = hit.Source.Decimals
			ethTokens[i].Address = hit.Source.Address
			ethTokens[i].Symbol = hit.Source.Symbol
			ethTokens[i].TotalSupply = hit.Source.TotalSupply
		}
	}

	return ethTokens, nil
}
