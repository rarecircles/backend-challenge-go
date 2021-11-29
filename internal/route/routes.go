package route

import (
	"net/http"

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
	r.Engine.Use(middleware.Logger())
	r.Engine.GET("/health-check", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	v1 := r.Engine.Group("/v1")
	v1.GET("/tokens", r.V1.Tokens)
}
