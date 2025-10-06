package repositories_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/FIAP-11SOAT/totem-de-pedidos/pkg/tests"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		panic(err)
	}

	productRepository := repositories.NewProductRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should failed with empty Product entity", func(t *testing.T) {
		_, err := productRepository.CreateProduct(context.Background(), &entity.Product{})
		assert.Error(t, err)
		assert.NotNil(t, err)
	})
}
