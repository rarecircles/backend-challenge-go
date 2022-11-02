package service

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	"go.uber.org/zap"
)

type SearchService interface {
	Search(ctx context.Context, token string) ([]eth.Token, error)
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

func (s *searchService) Search(ctx context.Context, token string) ([]eth.Token, error) {
	// TODO: search

	return nil, nil
}
