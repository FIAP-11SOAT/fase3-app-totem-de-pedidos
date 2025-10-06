package mapper

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/output"
)

func MapProductToOutput(product *entity.Product) output.ProductOutput {
	out := output.ProductOutput{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		Price:           product.Price,
		ImageURL:        product.ImageURL,
		PreparationTime: product.PreparationTime,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
		CategoryID:      product.CategoryID,
	}
	return out
}
