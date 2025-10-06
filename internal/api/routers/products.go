package routers

import (
	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/usecase"
	"github.com/labstack/echo/v4"
)

func ProductsRouter(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	r := repositories.NewProductRepository(dbConnection)
	u := usecase.NewProductUseCase(r)
	h := handlers.NewProductHandler(u)

	e.GET("/products", h.ListAllProducts)
	e.GET("/products/:id", h.FindProductByID)
	e.POST("/products", h.CreateProduct)
	e.PUT("/products/:id", h.UpdateProduct)
	e.DELETE("/products/:id", h.DeleteProduct)
}
