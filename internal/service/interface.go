package service

import "github.com/rarecircles/backend-challenge-go/internal/models"

type Interface interface {
	Search(string) ([]models.Token, error)
}

type IRedisRepo interface {
	Search(key string) ([]models.Token, error)
}
