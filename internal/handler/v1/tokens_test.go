package v1

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/jose-camilo/backend-challenge-go/internal/model"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/logging"
	searchEngineMock "github.com/jose-camilo/backend-challenge-go/internal/pkg/search_engine/mock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {

	tests := []struct {
		name  string
		id    string
		tokenDTO []model.TokenDTO
		error error
	}{
		{
			name:  "anySearchTerm",
			id:    "addressID",
			tokenDTO: []model.TokenDTO{
				model.TokenDTO{
					Name: "someName",
					Symbol: "someSymbol",
					Address: "0xsd89f7987ds9f8",
					Decimals: 2,
					TotalSupply: big.NewInt(9879878),
				},
			},
			error: nil,
		},
		{
			name:  "anySearchTerm",
			id:    "addressID",
			tokenDTO: []model.TokenDTO{},
			error: errors.New("test error"),
		},
	}
	{

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// ROUTE `/v1/tokens`
				req := httptest.NewRequest(http.MethodGet, "/v1/tokens?q="+test.name, nil)
				rec := httptest.NewRecorder()

				e := echo.New()
				echoCtx := e.NewContext(req, rec)

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				mockSearchEngine := searchEngineMock.NewMockSearchEngine(mockCtrl)
				mockSearchEngine.EXPECT().Search(gomock.Any(), test.name,"", "").Return(test.tokenDTO, test.error).Times(1)
				zLog := logging.MustCreateLoggerWithServiceName("challenge")

				handlers := &Handlers{
					SearchEngine: mockSearchEngine,
					Logger:       zLog,
				}

				err := handlers.Tokens(echoCtx)
				marshalled, _ := json.Marshal(apiResponse{
					"tokesn": test.tokenDTO,
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
