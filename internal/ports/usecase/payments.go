package usecase

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/output"
)

type Payments interface {
	CreatePayment(input input.CreatePaymentInput) (*output.CreatePaymentOutput, error)
	GetPaymentByID(paymentID string) (*entity.Payment, error)
	PaymentWebHook(*mercadopago.WebhookPayload) error
}
