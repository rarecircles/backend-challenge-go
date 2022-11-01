// Package tokengrp maintains the group of handlers for tokens
package tokengrp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/internal/model"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
}

func NewHandler(log *zap.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

type QueryTokensRequest struct {
	Query string `form:"q"`
}

type QueryTokensResponse struct {
	Tokens []model.Token `json:"tokens"`
}

func (h *Handler) QueryTokens(ctx *gin.Context) {
	var req QueryTokensRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("binding query: %w", err))
		return
	}

	h.log.Sugar().Infow("QueryTokens", "query", req)
	// TODO: query tokens

	var resp QueryTokensResponse
	resp.Tokens = []model.Token{}
	ctx.JSON(http.StatusOK, resp)
}
