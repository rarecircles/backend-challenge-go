package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	fixtures "github.com/rarecircles/backend-challenge-go/cmd/fixture"
	mocks "github.com/rarecircles/backend-challenge-go/cmd/mock"
	models "github.com/rarecircles/backend-challenge-go/cmd/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetTokens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		tokens        models.TokensDTO
		mockShop      func(mock *mocks.MockDaoInterface)
		expectedError error
	}{
		{
			name: "Happy path, tokens retrieved",
			tokens: models.TokensDTO{
				{
					Name:     "BitCoin",
					Symbol:   "BTC",
					Address:  "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
					Decimals: 18,
				},
			},
			mockShop: func(mock *mocks.MockDaoInterface) {
				mock.EXPECT().GetTokens("Bit").Return(fixtures.LoadTokensFixture("one_token"), nil)
			},
			expectedError: nil,
		},
		{
			name: "Error retriving tokens",
			mockShop: func(mock *mocks.MockDaoInterface) {
				mock.EXPECT().GetTokens("").Return(nil, errors.New("unable to connect to DB"))
			},
			expectedError: errors.New("unable to connect to DB"),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockDaoInterface(ctrl)
			test.mockShop(mockClient)
			tokens, err := GetTokens(mockClient, "Bit")
			if err != nil {
				assert.Equal(t, err, test.expectedError)
			} else {
				assert.Equal(t, tokens, test.tokens)
			}
		})
	}
}
