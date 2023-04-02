package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	fixtures "github.com/rarecircles/backend-challenge-go/cmd/fixture"
	mocks "github.com/rarecircles/backend-challenge-go/cmd/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetTokens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		mockShop func(mock *mocks.MockDaoInterface)
		status   int
	}{
		{
			name: "bad request requested",
			mockShop: func(mock *mocks.MockDaoInterface) {
				mock.EXPECT().GetTokens("").Return(nil, errors.New("unable to reach database"))
			},
			status: http.StatusBadRequest,
		},
		{
			name: "happy path, tokens retrived",
			mockShop: func(mock *mocks.MockDaoInterface) {
				mock.EXPECT().GetTokens("Rare").Return(fixtures.LoadTokensFixture("two_tokens"), nil)
			},
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", "/tokens?q=Year", errReader(0))
			if err != nil {
				t.Fatalf("Error creating a new request: %v", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockDaoInterface(ctrl)
			rr := httptest.NewRecorder()
			test.mockShop(mockClient)
			handler := http.HandlerFunc(GetTokens(mockClient))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
