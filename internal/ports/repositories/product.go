package repositories

import (
	"context"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type Product interface {
	ListProducts(ctx context.Context, input *input.ProductFilterInput) ([]*entity.Product, error)
	FindProductByID(ctx context.Context, id string) (*entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product) (int, error)
	GetCategoryByName(ctx context.Context, categoryName string) (*entity.ProductCategory, error)
	GetCategories(ctx context.Context) ([]*entity.ProductCategory, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	GetProductsByCategoryID(ctx context.Context, categoryID int) ([]*entity.Product, error)
}
