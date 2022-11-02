package token_grp_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	tokenGrp "github.com/rarecircles/backend-challenge-go/internal/api/handler/token_grp"
	"github.com/rarecircles/backend-challenge-go/internal/pkg/eth"
	mockSearch "github.com/rarecircles/backend-challenge-go/internal/service/search/mock"
	"github.com/rarecircles/backend-challenge-go/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryTokens(t *testing.T) {
	log := logger.MustCreateLoggerWithServiceName("TEST")

	tests := []struct {
		name       string
		query      string
		statusCode int
		expTokens  []eth.Token
		serviceErr error
		expErr     string
	}{
		{
			name:       "search service error",
			query:      "",
			statusCode: http.StatusInternalServerError,
			expTokens:  nil,
			serviceErr: errors.New("service error"),
			expErr:     "failed to search tokens:",
		},
		{
			name:       "success",
			query:      "rare",
			statusCode: http.StatusOK,
			expTokens: []eth.Token{
				{
					Name:        "RareCircles",
					Symbol:      "RCI",
					Address:     eth.MustNewAddress("dbf1344a0ff21bc098eb9ad4eef7de0f9722c02b"),
					Decimals:    18,
					TotalSupply: big.NewInt(1000000000),
				},
				{
					Name:        "Rareible",
					Symbol:      "RRI",
					Address:     eth.MustNewAddress("e9c8934ebd00bf73b0e961d1ad0794fb22837206"),
					Decimals:    9,
					TotalSupply: big.NewInt(100),
				},
			},
		},
		{
			name:       "success - empty array",
			query:      "ThereIsNoTokenHere",
			statusCode: http.StatusOK,
			expTokens:  []eth.Token{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tokens?q="+tt.query, nil)
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockSearchService := mockSearch.NewMockSearchService(mockCtrl)
			mockSearchService.EXPECT().SearchToken(gomock.Any(), tt.query).Return(tt.expTokens, tt.serviceErr).Times(1)

			h := tokenGrp.NewHandler(log, mockSearchService)
			h.QueryTokens(ctx)

			marshalled, err := json.Marshal(tokenGrp.QueryTokensResponse{
				Tokens: tokenGrp.ToTokenDTO(tt.expTokens),
			})
			require.NoError(t, err)

			if len(ctx.Errors) > 0 {
				assert.Equal(t, tt.statusCode, w.Code)
				assert.ErrorContains(t, ctx.Errors[0], tt.expErr)
			} else {
				assert.Equal(t, tt.statusCode, w.Code)
				assert.JSONEq(t, string(marshalled), w.Body.String())
			}
		})
	}
}
