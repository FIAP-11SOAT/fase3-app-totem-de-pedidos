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

type CustomerRepository struct {
	sqlClient *pgxpool.Pool
}

func NewCustomerRepository(database *dbadapter.DatabaseAdapter) repositories.Customer {
	return &CustomerRepository{
		sqlClient: database.Client,
	}
}

func (s *CustomerRepository) CreateCustomer(customer *entity.Customer) (*entity.Customer, error) {
	query := `INSERT INTO customers (name, email, tax_id) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.sqlClient.QueryRow(context.Background(), query, customer.Name, customer.Email, customer.TaxID).Scan(&customer.ID, &customer.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("error on inserting Customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerRepository) IdentifyCustomer(taxId *string) (*entity.Customer, error) {
	var customerResponse entity.Customer
	fmt.Println(taxId)
	err := s.sqlClient.QueryRow(context.Background(), "SELECT id, name, email, tax_id FROM customers WHERE tax_id = $1", taxId).Scan(&customerResponse.ID, &customerResponse.Name, &customerResponse.Email, &customerResponse.TaxID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error on query Customer: %w", err)
	}
	return &customerResponse, nil
}
