package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"

	"github.com/jackc/pgx/v5/pgxpool"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

const (
	GetOrderItemsQuery    = `SELECT id, product_id, quantity, price, created_at FROM order_items WHERE order_id = $1`
	GetOrderByIdQuery     = `SELECT id, status, total_amount, customer_id, order_date, created_at FROM orders WHERE id = $1`
	InsertOrderItemsQuery = `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
	InsertOrderQuery      = `INSERT INTO orders (status, total_amount, notification_attempts, customer_id) VALUES ($1, $2, $3, $4) RETURNING id`
	UpdateStatusQuery     = `UPDATE orders SET status = $1, updated_at = current_timestamp WHERE id = $2`
	ListOrdersBaseQuery   = `SELECT id, status, total_amount, customer_id, notification_attempts, order_date, created_at FROM orders`
	UpdatePaymentIDQuery  = `UPDATE orders SET payment_id = $1 WHERE id = $2`
)

type orderRepository struct {
	sqlClient *pgxpool.Pool
}

func NewOrderRepository(database *dbadapter.DatabaseAdapter) repositories.Order {
	return &orderRepository{
		sqlClient: database.Client,
	}
}

func (o *orderRepository) CreateOrder(input entity.Order) (int, error) {
	ctx := context.Background()

	var commited bool = false
	tx, err := o.sqlClient.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer func() {
		if !commited {
			tx.Rollback(ctx)
		}
	}()

	var orderID int
	err = tx.QueryRow(
		ctx, InsertOrderQuery,
		input.Status, input.TotalAmount, input.NotificationAttempts, input.CustomerID,
	).Scan(&orderID)
	if err != nil {
		return -1, fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range input.Items {
		_, err := tx.Exec(
			ctx, InsertOrderItemsQuery,
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			return -1, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	commited = true

	return orderID, nil
}

func (o *orderRepository) AddPayment(orderID int, paymentID int) error {
	_, err := o.sqlClient.Exec(context.Background(), UpdatePaymentIDQuery, paymentID, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order with payment ID: %w", err)
	}
	return nil
}

func (o *orderRepository) UpdateStatus(orderID int, status entity.OrderStatus) error {
	tag, err := o.sqlClient.Exec(context.Background(), UpdateStatusQuery, status, orderID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order with ID %d not found", orderID)
	}

	return nil
}

func (o *orderRepository) GetOrderByID(orderID int) (entity.Order, error) {
	row := o.sqlClient.QueryRow(context.Background(), GetOrderByIdQuery, orderID)

	var order entity.Order
	err := row.Scan(
		&order.ID,
		&order.Status,
		&order.TotalAmount,
		&order.CustomerID,
		&order.OrderDate,
		&order.CreatedAt,
	)
	if err != nil {
		return entity.Order{}, err
	}

	rows, err := o.sqlClient.Query(context.Background(), GetOrderItemsQuery, orderID)
	if err != nil {
		return entity.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem

		err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt)
		if err != nil {
			return entity.Order{}, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (o *orderRepository) ListOrders(filter input.OrderFilterInput) ([]entity.Order, error) {
	baseQuery := ListOrdersBaseQuery

	var where []string
	var args []interface{}
	argIndex := 1

	if filter.ID != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *filter.ID)
		argIndex++
	}

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.CustomerID != nil {
		where = append(where, fmt.Sprintf("customer_id = $%d", argIndex))
		args = append(args, *filter.CustomerID)
		argIndex++
	}

	if filter.NotificationAttempts != nil {
		where = append(where, fmt.Sprintf("notification_attempts = $%d", argIndex))
		args = append(args, *filter.NotificationAttempts)
		argIndex++
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	rows, err := o.sqlClient.Query(context.Background(), baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.TotalAmount,
			&order.CustomerID,
			&order.NotificationAttempts,
			&order.OrderDate,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
