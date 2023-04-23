package redisrepo

import (
	"errors"
	"github.com/RediSearch/redisearch-go/redisearch"
)

type redisClientMock struct {
}

func (r redisClientMock) Info() (*redisearch.IndexInfo, error) {
	return nil, nil
}

func (r redisClientMock) CreateIndex(*redisearch.Schema) (err error) {
	return nil
}

func (r redisClientMock) Index(docs ...redisearch.Document) error {
	return errors.New(docs[0].Properties["name"].(string) + docs[0].Properties["symbol"].(string))
}

func (r redisClientMock) Search(q *redisearch.Query) (docs []redisearch.Document, total int, err error) {
	return nil, 0, errors.New(q.Raw)
}
