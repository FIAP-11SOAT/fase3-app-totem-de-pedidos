package repositories

import "github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"

type Payments interface {
	GetPaymentByID(paymentID string) (*entity.Payment, error)
	CreatePayment(payment *entity.Payment) (*entity.Payment, error)
	UpdatePaymentStatus(paymentID string, status entity.PaymentStatus) error
	UpdatePaymentStatusWithOrderID(orderID int, status entity.PaymentStatus) error
}
