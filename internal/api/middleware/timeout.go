package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusRequestTimeout, gin.H{
		"mesage": "timeout",
	})
}

func TimeOut(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}
