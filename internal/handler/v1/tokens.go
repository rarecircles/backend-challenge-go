package v1

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/jose-camilo/backend-challenge-go/internal/model"
	"github.com/jose-camilo/backend-challenge-go/internal/pkg/search_engine"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	Logger *zap.Logger
	SearchEngine search_engine.SearchEngine
}
type apiResponse map[string]interface{}
type TokenResponseArray []model.TokenDTO

func (hv1 *Handlers) Tokens(ctx echo.Context) error {

	queryString := ctx.QueryParam("q")
	tokensSearchResponse, err := hv1.SearchEngine.Search(context.Background(), queryString, "","")
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse{
		"tokesn": tokensSearchResponse,
	})
}
