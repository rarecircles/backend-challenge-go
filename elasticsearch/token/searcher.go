package token

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/degarajesh/backend-challenge-go/model"
	v7 "github.com/elastic/go-elasticsearch/v7"
	v7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

type Searcher interface {
	Index(ctx context.Context, token model.Token) error
	Search(ctx context.Context, name string) ([]model.Token, error)
}

type ESSearcher struct {
	client *v7.Client
	input  chan model.Token
	zLog   *zap.Logger
	index  string
}

func NewEsSearcher(address string, input chan model.Token, logger *zap.Logger) (*ESSearcher, error) {
	config := v7.Config{
		Addresses: []string{
			address,
		},
	}
	esClient, err := v7.NewClient(config)
	if err != nil {
		return nil, err
	}

	logger.Info(v7.Version)

	return &ESSearcher{
		client: esClient,
		input:  input,
		zLog:   logger,
		index:  "token",
	}, nil
}

func (e *ESSearcher) Index(ctx context.Context, token model.Token) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(token)
	if err != nil {
		return err
	}

	req := v7api.IndexRequest{
		Index:      e.index,
		Body:       &buf,
		DocumentID: token.Symbol,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, e.client)
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

func (e *ESSearcher) Search(ctx context.Context, name string) ([]model.Token, error) {

	var query map[string]interface{}
	if name != "" {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"query_string": map[string]interface{}{
					"query":  name + "*",
					"fields": []string{"name"},
				},
			},
		}
	}

	req := v7api.SearchRequest{
		Index: []string{e.index},
	}
	var buf bytes.Buffer
	if query != nil {
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return nil, err
		}
		req.Body = &buf
	}

	resp, err := req.Do(ctx, e.client)
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
				Source model.Token `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		return nil, err
	}

	res := make([]model.Token, len(hits.Hits.Hits))
	for i, hit := range hits.Hits.Hits {
		res[i].Name = hit.Source.Name
		res[i].Decimals = hit.Source.Decimals
		res[i].Address = hit.Source.Address
		res[i].Symbol = hit.Source.Symbol
		res[i].TotalSupply = hit.Source.TotalSupply
	}
	return res, nil
}
