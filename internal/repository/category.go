package repositories

import (
	"context"
	"fmt"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	sqlClient *pgxpool.Pool
}

func NewCategoryRepository(database *dbadapter.DatabaseAdapter) repositories.Category {
	return &categoryRepository{
		sqlClient: database.Client,
	}
}

func (r *categoryRepository) GetCategories(ctx context.Context) ([]*entity.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM product_categories
	`

	rows, err := r.sqlClient.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}
	return categories, nil
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *entity.Category) (int, error) {
	query := `
		INSERT INTO product_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.sqlClient.QueryRow(ctx, query, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting category: %w", err)
	}
	return id, nil
}

func (r *categoryRepository) FindCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM product_categories
		WHERE id = $1
	`

	var category entity.Category
	err := r.sqlClient.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("error querying category: %w", err)
	}
	return &category, nil
}
func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	query := `
		UPDATE product_categories
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, description, created_at, updated_at
	`

	var updatedCategory entity.Category
	err := r.sqlClient.QueryRow(ctx, query, category.Name, category.Description, category.UpdatedAt, category.ID).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description, &updatedCategory.CreatedAt, &updatedCategory.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating category: %w", err)
	}
	return &updatedCategory, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, categoryID int) error {
	query := `
		DELETE FROM product_categories
		WHERE id = $1
	`

	_, err := r.sqlClient.Exec(ctx, query, categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	return nil
}
