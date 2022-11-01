package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck health check handler
func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
