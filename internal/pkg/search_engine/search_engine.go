package search_engine

//go:generate mockgen -destination=mock/mock_search_engine.go -package=mock github.com/jose-camilo/backend-challenge-go/internal/pkg/search_engine SearchEngine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jose-camilo/backend-challenge-go/internal/model"
	"go.uber.org/zap"
)

type SearchEngine interface {
	Index(ctx context.Context, token model.TokenDTO) error
	Search(ctx context.Context, name string, symbol string, address string) ([]model.TokenDTO, error)
}

type SearchEngineImpl struct {
	client *esv7.Client
	input  chan model.TokenDTO
	zLog   *zap.Logger
	index  string
}

func NewElasticsearchIngest(input chan model.TokenDTO, logger *zap.Logger) (SearchEngine, error) {
	config := esv7.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_HOSTS"),
		},
	}
	esClient, err := esv7.NewClient(config)
	if err != nil {
		return nil, err
	}

	logger.Info(esv7.Version)

	return &SearchEngineImpl{
		client: esClient,
		input:  input,
		zLog:   logger,
		index:  "token",
	}, nil
}

func (es *SearchEngineImpl) Index(ctx context.Context, token model.TokenDTO) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(token)
	if err != nil {
		return err
	}

	req := esv7api.IndexRequest{
		Index:      es.index,
		Body:       &buf,
		DocumentID: token.Symbol,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, es.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return err
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (es *SearchEngineImpl) Search(ctx context.Context, name string, symbol string, address string) ([]model.TokenDTO, error) {
	if name == "" && symbol == "" && address == "" {
		return nil, nil
	}

	should := make([]interface{}, 0, 3)

	if name != "" {
		should = append(should, map[string]interface{}{
			"fuzzy": map[string]interface{}{
				"name": name,
			},
		})
	}

	if symbol != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"symbol": symbol,
			},
		})
	}

	if address != "" {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"address": address,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	req := esv7api.SearchRequest{
		Index: []string{es.index},
		Body:  &buf,
	}
	es.zLog.Info(fmt.Sprintf("%+v", req.Body))

	resp, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, err
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source model.TokenDTO `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, err
	}

	res := make([]model.TokenDTO, len(hits.Hits.Hits))
	for i, hit := range hits.Hits.Hits {
		res[i].Name = hit.Source.Name
		res[i].Decimals = hit.Source.Decimals
		res[i].Address = hit.Source.Address
		res[i].Symbol = hit.Source.Symbol
		res[i].TotalSupply = hit.Source.TotalSupply
	}
	return res, nil
}
