// Package tokengrp maintains the group of handlers for tokens
package token_grp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	log           *zap.Logger
	searchService service.SearchService
}

func NewHandler(log *zap.Logger, searchService service.SearchService) *Handler {
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
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("binding query: %w", err))
		return
	}

	ethTokens, err := h.searchService.Search(ctx, req.Query)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to search tokens: %w", err))
		return
	}

	var resp QueryTokensResponse
	resp.Tokens = toTokenDTO(ethTokens)

	ctx.JSON(http.StatusOK, resp)
}
