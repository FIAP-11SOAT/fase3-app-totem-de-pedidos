package service

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/output"
)

type PaymentService interface {
	GeneratePaymentQrCode(input *input.CreatePaymentInput) (*output.CreatePaymentOutput, error)
	GetPaymentByID(paymentID string) (*output.GetPaymentOutput, error)
}
