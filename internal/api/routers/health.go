package routers

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/labstack/echo/v4"
)

func HealthRouter(e *echo.Echo) {
	h := handlers.NewHealthHandler()
	e.GET("/health", h.HealthCheck)
}
