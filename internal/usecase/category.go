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

type Category struct {
	categoryRepository repositories.Category
}

func NewCategoryCase(repository repositories.Category) usecase.Category {
	return &Category{categoryRepository: repository}
}

func (c *Category) GetCategories() ([]*entity.Category, error) {
	categories, err := c.categoryRepository.GetCategories(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting categories: %w", err)
	}
	return categories, nil
}

func (c *Category) CreateCategory(categoryDTO *input.CategoryInput) (*entity.Category, error) {
	categoryToCreate := &entity.Category{
		Name:        categoryDTO.Name,
		Description: categoryDTO.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	createdCategoryID, err := c.categoryRepository.CreateCategory(
		context.Background(),
		categoryToCreate,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating category: %w", err)
	}

	return &entity.Category{
		ID:          createdCategoryID,
		Name:        categoryDTO.Name,
		Description: categoryDTO.Description,
	}, nil
}

func (c *Category) UpdateCategory(categoryDTO *entity.Category) (*entity.Category, error) {
	categoryToUpdate := &entity.Category{
		ID:          categoryDTO.ID,
		Name:        categoryDTO.Name,
		Description: categoryDTO.Description,
		CreatedAt:   categoryDTO.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}

	updatedCategory, err := c.categoryRepository.UpdateCategory(context.Background(), categoryToUpdate)
	if err != nil {
		return nil, fmt.Errorf("error updating category: %w", err)
	}

	return updatedCategory, nil
}

func (c *Category) DeleteCategory(categoryID int) error {
	err := c.categoryRepository.DeleteCategory(context.Background(), categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	return nil
}

func (c *Category) FindCategoryByID(categoryID int) (*entity.Category, error) {
	category, err := c.categoryRepository.FindCategoryByID(context.Background(), categoryID)
	if err != nil {
		return nil, fmt.Errorf("error finding category: %w", err)
	}
	return category, nil
}
