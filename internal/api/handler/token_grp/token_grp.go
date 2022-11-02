// Package token_grp maintains the group of handlers for tokens
package token_grp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/service/search"
	"go.uber.org/zap"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrSomethingWrong = errors.New("something wrong")
)

type Handler struct {
	log           *zap.Logger
	searchService search.SearchService
}

func NewHandler(log *zap.Logger, searchService search.SearchService) *Handler {
	return &Handler{
		log:           log,
		searchService: searchService,
	}
}

type QueryTokensRequest struct {
	Query string `form:"q"`
}

type QueryTokensResponse struct {
	Tokens []TokenDTO `json:"tokens"`
}

func (h *Handler) QueryTokens(ctx *gin.Context) {
	var req QueryTokensRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(fmt.Errorf("binding query: %w", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": ErrInvalidInput,
		})
		return
	}

	ethTokens, err := h.searchService.SearchToken(ctx, req.Query)
	if err != nil {
		ctx.Error(fmt.Errorf("failed to search tokens: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrSomethingWrong,
		})
		return
	}

	var resp QueryTokensResponse
	resp.Tokens = ToTokenDTO(ethTokens)

	ctx.JSON(http.StatusOK, resp)
}
