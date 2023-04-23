package service

import (
	"errors"
	"github.com/rarecircles/backend-challenge-go/internal/models"
)

type redisRepoMock struct {
}

func (r redisRepoMock) Search(key string) ([]models.Token, error) {
	return nil, errors.New(key)
}
