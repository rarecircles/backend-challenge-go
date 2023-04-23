package redisrepo

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/rarecircles/backend-challenge-go/internal/models"
)

type IRedisRepo interface {
	Store(models.Token) error
	Search(key string) ([]models.Token, error)
	GetAllAddresses() (map[string]bool, error)
}

type IRedisClient interface {
	Info() (*redisearch.IndexInfo, error)
	CreateIndex(*redisearch.Schema) (err error)
	Index(docs ...redisearch.Document) error
	Search(q *redisearch.Query) (docs []redisearch.Document, total int, err error)
}
