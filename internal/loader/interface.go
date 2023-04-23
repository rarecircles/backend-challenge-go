package loader

import "github.com/rarecircles/backend-challenge-go/internal/models"

type ILoader interface {
	RunLoader()
}

type IRedisRepo interface {
	Store(models.Token) error
	GetAllAddresses() (map[string]bool, error)
}
