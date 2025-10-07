package routers

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/labstack/echo/v4"
)

func DocsRouter(e *echo.Echo) {
	h := handlers.NewDocsHandler()
	e.GET("/docs/:path", h.Docs)
}
