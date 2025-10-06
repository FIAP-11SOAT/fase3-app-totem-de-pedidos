package mock

import (
	"context"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type ProductMock struct {
	ListProductsFunc            func(context.Context, *input.ProductFilterInput) ([]*entity.Product, error)
	FindProductByIDFunc         func(context.Context, string) (*entity.Product, error)
	CreateProductFunc           func(context.Context, *entity.Product) (int, error)
	GetCategoryByNameFunc       func(context.Context, string) (*entity.ProductCategory, error)
	GetCategoriesFunc           func(context.Context) ([]*entity.ProductCategory, error)
	UpdateProductFunc           func(context.Context, *entity.Product) (*entity.Product, error)
	DeleteProductFunc           func(context.Context, string) error
	GetProductsByCategoryIDFunc func(context.Context, int) ([]*entity.Product, error)
}

func NewProductRepositoryMock() *ProductMock {
	return &ProductMock{
		ListProductsFunc:            func(context.Context, *input.ProductFilterInput) ([]*entity.Product, error) { return nil, nil },
		FindProductByIDFunc:         func(context.Context, string) (*entity.Product, error) { return nil, nil },
		CreateProductFunc:           func(context.Context, *entity.Product) (int, error) { return 0, nil },
		GetCategoryByNameFunc:       func(context.Context, string) (*entity.ProductCategory, error) { return nil, nil },
		GetCategoriesFunc:           func(context.Context) ([]*entity.ProductCategory, error) { return nil, nil },
		UpdateProductFunc:           func(context.Context, *entity.Product) (*entity.Product, error) { return nil, nil },
		DeleteProductFunc:           func(context.Context, string) error { return nil },
		GetProductsByCategoryIDFunc: func(context.Context, int) ([]*entity.Product, error) { return nil, nil },
	}
}

func (p *ProductMock) ListProducts(ctx context.Context, input *input.ProductFilterInput) ([]*entity.Product, error) {
	return p.ListProductsFunc(ctx, input)
}

func (p *ProductMock) FindProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return p.FindProductByIDFunc(ctx, id)
}

func (p *ProductMock) CreateProduct(ctx context.Context, product *entity.Product) (int, error) {
	return p.CreateProductFunc(ctx, product)
}

func (p *ProductMock) GetCategoryByName(ctx context.Context, categoryName string) (*entity.ProductCategory, error) {
	return p.GetCategoryByNameFunc(ctx, categoryName)
}

func (p *ProductMock) GetCategories(ctx context.Context) ([]*entity.ProductCategory, error) {
	return p.GetCategoriesFunc(ctx)
}

func (p *ProductMock) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	return p.UpdateProductFunc(ctx, product)
}

func (p *ProductMock) DeleteProduct(ctx context.Context, productID string) error {
	return p.DeleteProductFunc(ctx, productID)
}

func (p *ProductMock) GetProductsByCategoryID(ctx context.Context, categoryID int) ([]*entity.Product, error) {
	return p.GetProductsByCategoryIDFunc(ctx, categoryID)
}
