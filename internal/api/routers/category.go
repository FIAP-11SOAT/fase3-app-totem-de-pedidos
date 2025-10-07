package routers

import (
	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/middleware"
	"github.com/labstack/echo/v4"
)

func CategoryRouter(e *echo.Echo, dbConnection *dbadapter.DatabaseAdapter) {
	h := handlers.NewCategoryHandler(dbConnection)

	e.GET("/categories", h.ListAllCategories)
	e.GET("/categories/:id", h.FindCategoryByID)
	e.POST("/categories", h.CreateCategory, middleware.RequireAdminRole())
	e.PUT("/categories/:id", h.UpdateCategory, middleware.RequireAdminRole())
	e.DELETE("/categories/:id", h.DeleteCategory, middleware.RequireAdminRole())
}
