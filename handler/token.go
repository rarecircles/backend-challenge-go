package handler

import (
	"context"
	"github.com/degarajesh/backend-challenge-go/elasticsearch/token"
	"github.com/degarajesh/backend-challenge-go/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	Searcher token.Searcher
}

func (h *Handlers) GetTokens(ctx echo.Context) error {
	tokensSearchResponse, err := h.Searcher.Search(context.Background(), ctx.QueryParam("q"))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string][]model.Token{
		"tokens": tokensSearchResponse,
	})
}
