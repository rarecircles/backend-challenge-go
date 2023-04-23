package storagerepo

import (
	"errors"
	"github.com/RediSearch/redisearch-go/redisearch"
)

type redisClientMock struct{}

func (rcm redisClientMock) Info() (*redisearch.IndexInfo, error) {
	return nil, nil
}

func (rcm redisClientMock) CreateIndex(*redisearch.Schema) (err error) {
	return nil
}

func (rcm redisClientMock) Index(docs ...redisearch.Document) error {
	return errors.New(docs[0].Properties["name"].(string) + docs[0].Properties["symbol"].(string))
}

func (rcm redisClientMock) Search(q *redisearch.Query) (docs []redisearch.Document, total int, err error) {
	return nil, 0, errors.New(q.Raw)
}
