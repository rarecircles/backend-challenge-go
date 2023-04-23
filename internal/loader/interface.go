package loader

import "github.com/rarecircles/backend-challenge-go/internal/models"

type I interface {
	RunLoader()
}

type iStorageRepo interface {
	Store(models.Token) error
	GetAllAddresses() (map[string]bool, error)
}
