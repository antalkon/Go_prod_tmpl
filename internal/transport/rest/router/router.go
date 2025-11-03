package router

import (
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/handlers"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, ping *handlers.PingHandler) {
	api := e.Group("/api/v1")
	api.GET("/ping", ping.Ping)
}
