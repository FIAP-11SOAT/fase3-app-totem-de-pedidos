package repositories_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	repositories "github.com/FIAP-11SOAT/totem-de-pedidos/internal/repository"
	"github.com/FIAP-11SOAT/totem-de-pedidos/pkg/tests"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		assert.NoError(t, err)
	}

	orderRepository := repositories.NewOrderRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should create order with items", func(t *testing.T) {
		// create customer
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id) VALUES ('test_user', 'test_user@example.com', '12345678900')
		`)
		assert.NoError(t, err)

		// get id from customer
		var customerID int
		err = client.QueryRow(context.Background(), `SELECT id FROM customers WHERE email = 'test_user@example.com'`).Scan(&customerID)
		assert.NoError(t, err)

		// insert into product_categories
		_, err = client.Exec(context.Background(), `
			INSERT INTO product_categories (name, description) VALUES ('Food', 'Food truck')
		`)
		assert.NoError(t, err)

		// get categoryId
		var categoryID int
		err = client.QueryRow(context.Background(), `SELECT id FROM product_categories WHERE name = 'Food'`).Scan(&categoryID)
		assert.NoError(t, err)

		// insert into products
		_, err = client.Exec(context.Background(), `
			INSERT INTO products (name, description, price, preparation_time, category_id)
			VALUES ('X-Burger', 'Just some hamburguer', 25.00, 10, $1)
		`, categoryID)
		assert.NoError(t, err)

		// get productId
		var productID int
		err = client.QueryRow(context.Background(), `SELECT id FROM products WHERE name = 'X-Burger'`).Scan(&productID)
		assert.NoError(t, err)

		// order input
		order := entity.Order{
			CustomerID:           &customerID,
			Status:               "PENDING",
			TotalAmount:          50.00,
			NotificationAttempts: 0,
			Items: []entity.OrderItem{
				{
					ProductID: productID,
					Quantity:  2,
					Price:     25.00,
				},
			},
			ID: customerID,
		}

		// create order
		orderID, err := orderRepository.CreateOrder(order)
		assert.NoError(t, err)
		assert.True(t, orderID > 0)
	})

	t.Run("should fail when product does not exist", func(t *testing.T) {
		// create customer
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id) VALUES ('error_user', 'error_user@example.com', '00000000000')
		`)
		assert.NoError(t, err)

		// get customerId
		var customerID int
		err = client.QueryRow(context.Background(), `SELECT id FROM customers WHERE email = 'error_user@example.com'`).Scan(&customerID)
		assert.NoError(t, err)

		// order input
		order := entity.Order{
			CustomerID:  &customerID,
			Status:      "PENDING",
			TotalAmount: 30.00,
			Items: []entity.OrderItem{
				{
					ProductID: 9999, // non-existent id
					Quantity:  1,
					Price:     30.00,
				},
			},
		}

		orderId, err := orderRepository.CreateOrder(order)
		assert.Error(t, err)
		assert.Equal(t, -1, orderId)
	})

}

func TestUpdateStatus(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		assert.NoError(t, err)
	}

	repo := repositories.NewOrderRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should update order status successfully", func(t *testing.T) {
		// cria customer
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id) VALUES ('status_user', 'status@example.com', '11122233344')
		`)
		assert.NoError(t, err)

		var customerID int
		err = client.QueryRow(context.Background(), `SELECT id FROM customers WHERE email = 'status@example.com'`).Scan(&customerID)
		assert.NoError(t, err)

		// cria order
		var orderID int
		err = client.QueryRow(context.Background(), `
			INSERT INTO orders (status, total_amount, notification_attempts, customer_id)
			VALUES ('PENDING', 100.00, 0, $1)
			RETURNING id
		`, customerID).Scan(&orderID)
		assert.NoError(t, err)

		// atualiza status
		err = repo.UpdateStatus(orderID, "COMPLETED")
		assert.NoError(t, err)

		// verifica no banco
		var updatedStatus string
		err = client.QueryRow(context.Background(), `SELECT status FROM orders WHERE id = $1`, orderID).Scan(&updatedStatus)
		assert.NoError(t, err)
		assert.Equal(t, "COMPLETED", updatedStatus)
	})

	t.Run("should return error if order id does not exist", func(t *testing.T) {
		err := repo.UpdateStatus(99999, "CANCELLED")
		assert.Error(t, err)
	})
}

func TestGetOrderByID(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		assert.NoError(t, err)
	}

	repo := repositories.NewOrderRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should return order and items", func(t *testing.T) {
		// create customer
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id)
			VALUES ('order_user', 'order@example.com', '12345678909')
		`)
		assert.NoError(t, err)

		// get customerId
		var customerID int
		err = client.QueryRow(context.Background(), `SELECT id FROM customers WHERE email = 'order@example.com'`).Scan(&customerID)
		assert.NoError(t, err)

		// create product category
		_, err = client.Exec(context.Background(), `
			INSERT INTO product_categories (name, description) VALUES ('TestCat', 'Desc')
		`)
		assert.NoError(t, err)

		// get category id
		var categoryID int
		err = client.QueryRow(context.Background(), `SELECT id FROM product_categories WHERE name = 'TestCat'`).Scan(&categoryID)
		assert.NoError(t, err)

		// insert product
		_, err = client.Exec(context.Background(), `
			INSERT INTO products (name, description, price, preparation_time, category_id)
			VALUES ('TestBurger', 'Delicious', 15.50, 10, $1)
		`, categoryID)
		assert.NoError(t, err)

		// get product id
		var productID int
		err = client.QueryRow(context.Background(), `SELECT id FROM products WHERE name = 'TestBurger'`).Scan(&productID)
		assert.NoError(t, err)

		// create order
		orderId, err := repo.CreateOrder(entity.Order{
			NotificationAttempts: 0,
			Status:               "PENDING",
			TotalAmount:          31.00,
			CustomerID:           &customerID,
		})
		assert.NoError(t, err)

		// create item
		_, err = client.Exec(context.Background(), `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, 2, 15.50)
		`, orderId, productID)
		assert.NoError(t, err)

		order, err := repo.GetOrderByID(orderId)
		assert.NoError(t, err)
		assert.Equal(t, orderId, order.ID)
		assert.Equal(t, "PENDING", string(order.Status))
		assert.Equal(t, &customerID, order.CustomerID)
		assert.Len(t, order.Items, 1)
		assert.Equal(t, productID, order.Items[0].ProductID)
		assert.Equal(t, 2, order.Items[0].Quantity)
	})
}

func TestListOrders(t *testing.T) {
	connStr := tests.CreatePostgresDataBase(t)

	client, err := pgxpool.NewWithConfig(context.Background(), dbadapter.Config(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}

	repo := repositories.NewOrderRepository(&dbadapter.DatabaseAdapter{Client: client})

	t.Run("should list orders with and without filters", func(t *testing.T) {
		// create customer
		_, err := client.Exec(context.Background(), `
			INSERT INTO customers (name, email, tax_id)
			VALUES ('filter_user', 'filter@example.com', '12345678999')
		`)
		assert.NoError(t, err)

		// get customer id
		var customerId int
		err = client.QueryRow(context.Background(), `SELECT id FROM customers WHERE email = 'filter@example.com'`).Scan(&customerId)
		assert.NoError(t, err)

		// insert orders
		for i := 1; i <= 3; i++ {
			status := "PENDING"
			if i == 2 {
				status = "COMPLETED"
			}

			_, err = client.Exec(context.Background(), `
				INSERT INTO orders (status, total_amount, customer_id, notification_attempts)
				VALUES ($1, $2, $3, $4)
			`, status, float64(i*10), customerId, i-1)
			assert.NoError(t, err)
		}

		orders, err := repo.ListOrders(input.OrderFilterInput{})
		assert.NoError(t, err)
		assert.Len(t, orders, 3)

		orders, err = repo.ListOrders(input.OrderFilterInput{Status: "COMPLETED"})
		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, "COMPLETED", string(orders[0].Status))

		orders, err = repo.ListOrders(input.OrderFilterInput{CustomerID: &customerId})
		assert.NoError(t, err)
		assert.Len(t, orders, 3)

		attempts := 0
		orders, err = repo.ListOrders(input.OrderFilterInput{NotificationAttempts: &attempts})
		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, 0, orders[0].NotificationAttempts)

		invalidID := 9999
		orders, err = repo.ListOrders(input.OrderFilterInput{ID: &invalidID})
		assert.NoError(t, err)
		assert.Len(t, orders, 0)
	})
}
