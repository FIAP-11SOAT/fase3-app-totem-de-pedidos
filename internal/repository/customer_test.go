package repositories_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/FIAP-11SOAT/totem-de-pedidos/pkg/tests"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepository(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		assert.NoError(t, err)
	}

	repo := repositories.NewCustomerRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should create customer", func(t *testing.T) {
		customer := &entity.Customer{
			Name:  "Test User",
			Email: "testuser@example.com",
			TaxID: "11122233344",
		}

		result, err := repo.CreateCustomer(customer)
		assert.NoError(t, err)
		assert.NotNil(t, result.ID)
		assert.NotZero(t, result.CreatedAt)

		// Confirm in DB
		var dbID int
		var dbEmail string
		err = client.QueryRow(context.Background(), `SELECT id, email FROM customers WHERE tax_id = $1`, customer.TaxID).
			Scan(&dbID, &dbEmail)
		assert.NoError(t, err)
		assert.Equal(t, customer.Email, dbEmail)
	})

	t.Run("should identify existing customer", func(t *testing.T) {
		taxID := "55566677788"
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id)
			VALUES ('Jane Doe', 'jane@example.com', $1)
		`, taxID)
		assert.NoError(t, err)

		foundCustomer, err := repo.IdentifyCustomer(&taxID)
		assert.NoError(t, err)
		assert.NotNil(t, foundCustomer)
		assert.Equal(t, "jane@example.com", foundCustomer.Email)
	})

	t.Run("should return nil if customer not found", func(t *testing.T) {
		taxID := "00000000000"
		customer, err := repo.IdentifyCustomer(&taxID)
		assert.NoError(t, err)
		assert.Nil(t, customer)
	})
}
