package repositories

import (
	"context"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
)

type Category interface {
	GetCategories(ctx context.Context) ([]*entity.Category, error)
	CreateCategory(ctx context.Context, category *entity.Category) (int, error)
	FindCategoryByID(ctx context.Context, id int) (*entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	DeleteCategory(ctx context.Context, categoryID int) error
}
