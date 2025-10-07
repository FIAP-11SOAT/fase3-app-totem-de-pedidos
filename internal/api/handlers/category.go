package handlers

import (
	"net/http"
	"strconv"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	service "github.com/FIAP-11SOAT/totem-de-pedidos/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService usecase.Category
}

func NewCategoryHandler(dbConnection *dbadapter.DatabaseAdapter) *CategoryHandler {
	categoryRepository := repositories.NewCategoryRepository(dbConnection)

	return &CategoryHandler{
		categoryService: service.NewCategoryCase(categoryRepository),
	}
}

func (c *CategoryHandler) ListAllCategories(ctx echo.Context) error {
	categories, err := c.categoryService.GetCategories()
	if err != nil {
		ctx.Logger().Error("Error getting categories", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if len(categories) == 0 {
		return ctx.JSON(http.StatusNotFound, make([]string, 0))
	}

	return ctx.JSON(http.StatusOK, categories)
}

func (c *CategoryHandler) FindCategoryByID(ctx echo.Context) error {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Logger().Error("Error converting category ID", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	category, err := c.categoryService.FindCategoryByID(categoryID)
	if err != nil {
		ctx.Logger().Error("Error finding category", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if category == nil {
		return ctx.JSON(http.StatusNotFound, make([]string, 0))
	}

	return ctx.JSON(http.StatusOK, category)
}

func (c *CategoryHandler) CreateCategory(ctx echo.Context) error {
	categoryInput := new(input.CategoryInput)
	if err := ctx.Bind(categoryInput); err != nil {
		ctx.Logger().Error("Error binding category input", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := categoryInput.Validate(); err != nil {
		ctx.Logger().Error("Error validating category input", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	category, err := c.categoryService.CreateCategory(categoryInput)
	if err != nil {
		ctx.Logger().Error("Error creating category", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryHandler) UpdateCategory(ctx echo.Context) error {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Logger().Error("Error converting category ID", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	categoryInput := new(entity.Category)
	if err := ctx.Bind(categoryInput); err != nil {
		ctx.Logger().Error("Error binding category input", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	categoryInput.ID = categoryID

	category, err := c.categoryService.UpdateCategory(categoryInput)
	if err != nil {
		ctx.Logger().Error("Error updating category", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, category)
}

func (c *CategoryHandler) DeleteCategory(ctx echo.Context) error {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Logger().Error("Error converting category ID", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	err = c.categoryService.DeleteCategory(categoryID)
	if err != nil {
		ctx.Logger().Error("Error deleting category", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}
