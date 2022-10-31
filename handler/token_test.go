package handler

import (
	"encoding/json"
	"errors"
	"github.com/degarajesh/backend-challenge-go/elasticsearch/mock"
	"github.com/degarajesh/backend-challenge-go/model"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

type apiResponse map[string]interface{}

func TestGetTokensHandler(t *testing.T) {

	tests := []struct {
		name     string
		id       string
		tokenDTO []model.Token
		error    error
	}{
		{
			name: "searchTerm",
			id:   "addressID",
			tokenDTO: []model.Token{
				model.Token{
					Name:        "RareCircles",
					Symbol:      "RCI",
					Address:     "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
					Decimals:    18,
					TotalSupply: big.NewInt(1000000000),
				},
			},
			error: nil,
		},
		{
			name:     "error",
			id:       "addressID",
			tokenDTO: []model.Token{},
			error:    errors.New("test error"),
		},
	}
	{

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/tokens?q="+test.name, nil)
				rec := httptest.NewRecorder()

				e := echo.New()
				echoCtx := e.NewContext(req, rec)

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				mockSearchEngine := mock.NewMockSearchEngine(mockCtrl)
				mockSearchEngine.EXPECT().Search(gomock.Any(), test.name).Return(test.tokenDTO, test.error).Times(1)

				handlers := &Handlers{
					Searcher: mockSearchEngine,
				}
				err := handlers.GetTokens(echoCtx)
				marshalled, _ := json.Marshal(apiResponse{
					"tokens": test.tokenDTO,
				})

				if test.error == nil {
					assert.JSONEqf(t, string(marshalled), rec.Body.String(), "")
					assert.Equal(t, http.StatusOK, rec.Code)
				} else {
					assert.EqualValues(t, test.error.Error(), err.Error())
				}
			})

		}
	}
}

func TestGetAllTokensHandler(t *testing.T) {

	tests := []struct {
		name     string
		id       string
		tokenDTO []model.Token
		error    error
	}{
		{
			name: "searchTerm",
			id:   "addressID",
			tokenDTO: []model.Token{
				model.Token{
					Name:        "RareCircles",
					Symbol:      "RCI",
					Address:     "0xdbf1344a0ff21bc098eb9ad4eef7de0f9722c02b",
					Decimals:    18,
					TotalSupply: big.NewInt(1000000000),
				},
			},
			error: nil,
		},
		{
			name:     "error",
			id:       "addressID",
			tokenDTO: []model.Token{},
			error:    errors.New("test error"),
		},
	}
	{

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/tokens", nil)
				rec := httptest.NewRecorder()

				e := echo.New()
				echoCtx := e.NewContext(req, rec)

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				mockSearchEngine := mock.NewMockSearchEngine(mockCtrl)
				mockSearchEngine.EXPECT().Search(gomock.Any(), "").Return(test.tokenDTO, test.error).Times(1)

				handlers := &Handlers{
					Searcher: mockSearchEngine,
				}
				err := handlers.GetTokens(echoCtx)
				marshalled, _ := json.Marshal(apiResponse{
					"tokens": test.tokenDTO,
				})

				if test.error == nil {
					assert.JSONEqf(t, string(marshalled), rec.Body.String(), "")
					assert.Equal(t, http.StatusOK, rec.Code)
				} else {
					assert.EqualValues(t, test.error.Error(), err.Error())
				}
			})

		}
	}
}
