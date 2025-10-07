package repositories

import (
	"context"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
)

const (
	CreatePaymentQuery                  = `INSERT INTO payments (amount, status, provider) VALUES ($1, $2, $3) RETURNING id`
	GetPaymentByIDQuery                 = `SELECT id, amount, payment_date, status, provider FROM payments WHERE id = $1`
	UpdatePaymentStatusQuery            = `UPDATE payments SET status = $1 WHERE id = $2`
	UpdatePaymentStatusWithOrderIDQuery = `UPDATE payments SET status = $1 WHERE id = ( SELECT payment_id FROM orders WHERE id = $2 );`
)

type paymentsRepository struct {
	sqlClient *pgxpool.Pool
}

func NewPaymentsRepository(database *dbadapter.DatabaseAdapter) repositories.Payments {
	return &paymentsRepository{
		sqlClient: database.Client,
	}
}

func (r *paymentsRepository) UpdatePaymentStatusWithOrderID(orderID int, status entity.PaymentStatus) error {
	_, err := r.sqlClient.Exec(context.Background(), UpdatePaymentStatusWithOrderIDQuery, status, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentsRepository) UpdatePaymentStatus(paymentID string, status entity.PaymentStatus) error {
	_, err := r.sqlClient.Exec(context.Background(), UpdatePaymentStatusQuery, status, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentsRepository) CreatePayment(payment *entity.Payment) (*entity.Payment, error) {
	err := r.sqlClient.
		QueryRow(context.Background(), CreatePaymentQuery,
			payment.Amount,
			payment.Status,
			payment.Provider,
		).Scan(&payment.ID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentsRepository) GetPaymentByID(paymentID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.sqlClient.
		QueryRow(context.Background(), GetPaymentByIDQuery, paymentID).
		Scan(
			&payment.ID,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.Status,
			&payment.Provider,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &payment, nil
		}
		return &payment, err
	}
	return &payment, nil
}
