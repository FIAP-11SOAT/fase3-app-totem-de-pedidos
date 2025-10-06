package routers

import (
	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/usecase"
	"github.com/labstack/echo/v4"
)

func CustomerRouter(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	customerRepository := repositories.NewCustomerRepository(dbConnection)
	customerUseCase := usecase.NewCustomerUseCase(customerRepository)
	customerHandler := handlers.NewCustomerHandler(customerUseCase)

	e.POST("/customer", customerHandler.CreateCustomer)
	e.GET("/customer", customerHandler.IdentifyCustomer)
}
