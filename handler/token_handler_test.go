package handler_test

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	appContext "github.com/rarecircles/backend-challenge-go/app_context.go"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/router"
	"github.com/rarecircles/backend-challenge-go/service"
	"github.com/rarecircles/backend-challenge-go/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type tokenHandlerTestSuite struct {
	suite.Suite
	router  http.Handler
	service *mocks.MockTokenService
}

func (suite *tokenHandlerTestSuite) SetupSuite() {
	rpcUrl := "https://rpc.com/shree/" // Ideally this should be retreived from test config
	appContext.Init(rpcUrl)
	suite.service = &mocks.MockTokenService{}
	dependencies := service.ServerDependencies{
		TokenService: suite.service,
	}
	suite.router = router.InitRouter(router.Options{
		Dependencies: &dependencies,
	})
}

func (suite *tokenHandlerTestSuite) TearDownSuite() {
	fmt.Println("TEARING DOWN Token Handler Suite")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(tokenHandlerTestSuite))
}

func (suite *tokenHandlerTestSuite) TestGetTokenInfo() {
	testCases := []struct {
		name               string
		queryParam         string
		mockedData         []eth.Token
		mockedErr          error
		expectedStatusCode int
	}{
		{
			name:       "test with valid data",
			queryParam: "rare",
			mockedData: []eth.Token{
				{
					Name:        "rarecircles",
					Symbol:      "RCI",
					Address:     []byte("0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02"),
					Decimals:    19,
					TotalSupply: big.NewInt(1000000000),
				},
			},
			expectedStatusCode: 200,
		},
		{
			name: "test with no query param",
			mockedData: []eth.Token{
				{
					Name:        "rarecircles",
					Symbol:      "RCI",
					Address:     []byte("0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02"),
					Decimals:    18,
					TotalSupply: big.NewInt(1000000000),
				},
			},
			expectedStatusCode: 400,
		},
		{
			name:               "test with service error",
			queryParam:         "rare",
			mockedErr:          errors.New(""),
			expectedStatusCode: 500,
		},
	}

	for _, t := range testCases {
		fmt.Printf("Running test: %v", t.name)
		suite.service.On("GetTokensInfo", mock.Anything).Return(t.mockedData, t.mockedErr).Once()

		endpoint := fmt.Sprintf("/tokens?q=%v", t.queryParam)
		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		response := httptest.NewRecorder()
		suite.router.ServeHTTP(response, req)
		suite.Equal(t.expectedStatusCode, response.Code)
		suite.service.Mock.ExpectedCalls = nil
	}
}
