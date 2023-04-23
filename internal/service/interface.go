package service

import "github.com/rarecircles/backend-challenge-go/internal/models"

type I interface {
	Search(string) ([]models.Token, error)
}

type iStorageRepo interface {
	Search(key string) ([]models.Token, error)
}
