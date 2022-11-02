package middleware

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func keyFunc(ctx *gin.Context) string {
	return ctx.ClientIP()
}

func errorHandler(ctx *gin.Context, info ratelimit.Info) {
	ctx.JSON(http.StatusTooManyRequests, gin.H{
		"message": "Too many requests. Try again in " + time.Until(info.ResetTime).String(),
	})
}

func NewRateLimiter(limit uint, rate time.Duration) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  rate,
		Limit: limit,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return rateLimiter
}
