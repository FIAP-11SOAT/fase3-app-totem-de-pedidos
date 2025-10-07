package usecase

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type Category interface {
	GetCategories() ([]*entity.Category, error)
	CreateCategory(categoryDTO *input.CategoryInput) (*entity.Category, error)
	UpdateCategory(categoryDTO *entity.Category) (*entity.Category, error)
	DeleteCategory(categoryID int) error
	FindCategoryByID(categoryID int) (*entity.Category, error)
}
