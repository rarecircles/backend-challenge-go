package mocks

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GetTokensInfo(tokenTitle string) ([]eth.Token, error) {
	args := m.Mock.Called(tokenTitle)
	return args.Get(0).([]eth.Token), args.Error(1)
}
