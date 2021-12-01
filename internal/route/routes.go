package route

import (
	"time"

	"github.com/jose-camilo/backend-challenge-go/internal/handler"
	v1Handler "github.com/jose-camilo/backend-challenge-go/internal/handler/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	Engine         *echo.Echo
	V1             v1Handler.Handlers
	CommonHandlers *handler.Handlers
	AddressChannel chan string
}

func (r *Router) Init() {
	r.Engine.Use(middleware.RequestID())
	r.Engine.Use(middleware.Timeout())
	r.Engine.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	r.Engine.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))

	r.Engine.Use(middleware.Logger())
	r.Engine.GET("/health-check", r.CommonHandlers.HealthCheck)

	v1 := r.Engine.Group("/v1")
	v1.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	v1.GET("/tokens", r.V1.Tokens)
	v1.GET("/symbols", r.V1.Symbols)
	v1.GET("/addresses", r.V1.Addresses)
}
