package service

import (
	"errors"
	"github.com/rarecircles/backend-challenge-go/internal/models"
)

type storageRepoMock struct{}

func (rrm storageRepoMock) Search(key string) ([]models.Token, error) {
	return nil, errors.New(key)
}
