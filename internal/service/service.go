package service

import (
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
)

type S struct {
	storageRepo iStorageRepo
	logger      *zap.Logger
}

func New(l *zap.Logger, storageRepo iStorageRepo) S {
	s := S{logger: l, storageRepo: storageRepo}
	return s
}

func (s S) Search(key string) ([]models.Token, error) {
	tokens, err := s.storageRepo.Search(key)
	return tokens, err
}
