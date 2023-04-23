package service

import (
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
)

type S struct {
	r IRedisRepo
	l *zap.Logger
}

func New(l *zap.Logger, repo IRedisRepo) S {
	s := S{l: l, r: repo}
	return s
}

func (s S) Search(key string) ([]models.Token, error) {
	tokens, err := s.r.Search(key)
	return tokens, err
}
