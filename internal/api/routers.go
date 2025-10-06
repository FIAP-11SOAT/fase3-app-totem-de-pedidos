package api

import (
	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/routers"
	"github.com/labstack/echo/v4"
)

func Routers(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	routers.DocsRouter(e)
	routers.HealthRouter(e)
	routers.ProductsRouter(e, dbConnection)
	routers.CustomerRouter(e, dbConnection)
	routers.CategoryRouter(e, dbConnection)
	routers.PaymentsRouter(e, dbConnection)
	routers.OrdersRouter(e, dbConnection)
}
