package usecase

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type Product interface {
	GetProductByID(id string) (*entity.Product, error)
	GetProducts(*input.ProductFilterInput) ([]*entity.Product, error)
	CreateProduct(productDTO *input.ProductInput) (*entity.Product, error)
	UpdateProduct(id string, productDTO *input.ProductInput) (*entity.Product, error)
	DeleteProduct(productID string) error
}
