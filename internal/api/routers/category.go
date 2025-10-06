package routers

import (
	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/labstack/echo/v4"
)

func CategoryRouter(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	h := handlers.NewCategoryHandler(dbConnection)

	e.GET("/categories", h.ListAllCategories)
	e.GET("/categories/:id", h.FindCategoryByID)
	e.POST("/categories", h.CreateCategory)
	e.PUT("/categories/:id", h.UpdateCategory)
	e.DELETE("/categories/:id", h.DeleteCategory)
}
