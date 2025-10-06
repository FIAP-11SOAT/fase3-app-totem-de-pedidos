package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
)

type Product struct {
	productRepository repositories.Product
}

func NewProductUseCase(repository repositories.Product) usecase.Product {
	return &Product{productRepository: repository}
}

func (p *Product) GetProductByID(id string) (*entity.Product, error) {
	product, err := p.productRepository.FindProductByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("error fetching product by id: %w", err)
	}
	return product, nil
}

func (p *Product) GetProducts(input *input.ProductFilterInput) ([]*entity.Product, error) {
	products, err := p.productRepository.ListProducts(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *Product) CreateProduct(productInput *input.ProductInput) (*entity.Product, error) {
	productToCreate := &entity.Product{
		Name:            productInput.Name,
		Description:     productInput.Description,
		Price:           productInput.Price,
		ImageURL:        productInput.ImageURL,
		PreparationTime: 0, // ajuste se necessário
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		CategoryID:      productInput.CategoryID,
	}

	createdProductID, err := p.productRepository.CreateProduct(
		context.Background(),
		productToCreate,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	productToCreate.ID = createdProductID
	return productToCreate, nil
}

func (p *Product) UpdateProduct(id string, productInput *input.ProductInput) (*entity.Product, error) {
	existing, err := p.productRepository.FindProductByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("error fetching product: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("product not found")
	}

	existing.Name = productInput.Name
	existing.Description = productInput.Description
	existing.Price = productInput.Price
	existing.ImageURL = productInput.ImageURL
	existing.CategoryID = productInput.CategoryID
	existing.UpdatedAt = time.Now().UTC()

	updated, err := p.productRepository.UpdateProduct(context.Background(), existing)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}
	return updated, nil
}

func (p *Product) DeleteProduct(productID string) error {
	err := p.productRepository.DeleteProduct(context.Background(), productID)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}
	return nil
}
