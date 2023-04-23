package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type H struct {
	svc    service.I
	logger *zap.Logger
}

func New(l *zap.Logger, svc service.I) H {
	return H{svc: svc, logger: l}
}

func (h H) Up(c *gin.Context) {
	c.JSON(http.StatusOK, "Tokens service is up")
}

func (h H) SearchTokens(c *gin.Context) {
	params := c.Request.URL.Query()
	q := params.Get("q")

	tokens, err := h.svc.Search(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tokens == nil {
		c.JSON(http.StatusOK, []eth.Token{})
		return
	}
	c.JSON(http.StatusOK, tokens)
}
