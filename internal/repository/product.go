package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
)

type productRepository struct {
	sqlClient *pgxpool.Pool
}

func NewProductRepository(database *dbadapter.DatabaseAdapter) repositories.Product {
	return &productRepository{
		sqlClient: database.Client,
	}
}

func (p *productRepository) ListProducts(ctx context.Context, input *input.ProductFilterInput) ([]*entity.Product, error) {
	query := `
        SELECT p.id, p.name, p.description, p.price, p.image_url, p.preparation_time, p.created_at, p.updated_at, p.category_id
        FROM products p
        JOIN product_categories pc ON p.category_id = pc.id
    `
	var args []interface{}
	where := ""

	if input.Name != "" {
		args = append(args, "%"+input.Name+"%")
		where += fmt.Sprintf("p.name ILIKE $%d", len(args))
	}
	if input.CategoryName != "" {
		if where != "" {
			where += " AND "
		}
		args = append(args, "%"+input.CategoryName+"%")
		where += fmt.Sprintf("pc.name ILIKE $%d", len(args))
	}
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := p.sqlClient.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.PreparationTime,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CategoryID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

func (p *productRepository) FindProductByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, preparation_time, created_at, updated_at, category_id
		FROM products
		WHERE id = $1
	`

	var product entity.Product
	err := p.sqlClient.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageURL,
		&product.PreparationTime,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CategoryID,
	)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &product, nil
}

func (p *productRepository) CreateProduct(ctx context.Context, product *entity.Product) (int, error) {
	var createdProductID int

	query := `
		INSERT INTO products (name, description, price, image_url, preparation_time, created_at, updated_at, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := p.sqlClient.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageURL,
		product.PreparationTime,
		product.CreatedAt,
		product.UpdatedAt,
		product.CategoryID,
	).Scan(&createdProductID)
	if err != nil {
		return 0, err
	}

	return createdProductID, nil
}

func (p *productRepository) GetProductsByCategoryID(ctx context.Context, categoryID int) ([]*entity.Product, error) {
	query := `
        SELECT id, name, description, price, image_url, preparation_time, created_at, updated_at, category_id
        FROM products
        WHERE category_id = $1
    `

	rows, err := p.sqlClient.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error querying products by category ID: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.PreparationTime,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CategoryID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over product rows: %w", err)
	}

	return products, nil
}

func (p *productRepository) GetCategoryByName(ctx context.Context, categoryName string) (*entity.ProductCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM product_categories
		WHERE name = $1
	`

	var category entity.ProductCategory
	err := p.sqlClient.QueryRow(ctx, query, categoryName).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	query := `
		UPDATE products
		SET name = $1,
			description = $2,
			price = $3,
			image_url = $4,
			preparation_time = $5,
			updated_at = $6,
			category_id = $7
		WHERE id = $8
		RETURNING id, name, description, price, image_url, preparation_time, created_at, updated_at, category_id
	`

	var updatedProduct entity.Product
	err := p.sqlClient.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageURL,
		product.PreparationTime,
		product.UpdatedAt,
		product.CategoryID,
		product.ID,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Description,
		&updatedProduct.Price,
		&updatedProduct.ImageURL,
		&updatedProduct.PreparationTime,
		&updatedProduct.CreatedAt,
		&updatedProduct.UpdatedAt,
		&updatedProduct.CategoryID,
	)

	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, productID string) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	_, err := p.sqlClient.Exec(ctx, query, productID)
	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetCategories(ctx context.Context) ([]*entity.ProductCategory, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM product_categories
	`

	rows, err := p.sqlClient.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.ProductCategory
	for rows.Next() {
		var category entity.ProductCategory
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return categories, nil
}
