package usecase

import (
	"errors"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/output"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/service"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
)

type Payments struct {
	paymentRepository repositories.Payments
	paymentService    service.PaymentService
	orderRepository   repositories.Order
}

func NewPaymentsUseCase(
	paymentRepository repositories.Payments,
	paymentService service.PaymentService,
	orderRepository repositories.Order,
) usecase.Payments {
	return &Payments{
		paymentRepository: paymentRepository,
		paymentService:    paymentService,
		orderRepository:   orderRepository,
	}
}

func (p *Payments) GetPaymentByID(paymentID string) (*entity.Payment, error) {
	payment, err := p.paymentRepository.GetPaymentByID(paymentID)
	if err != nil {
		return payment, err
	}
	return payment, nil
}

func (p *Payments) PaymentWebHook(payload *mercadopago.WebhookPayload) error {
	paymentID := payload.Data.ID
	result, err := p.paymentService.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}

	if result.Status != "approved" {
		return errors.New("payment not approved")
	}

	err = p.orderRepository.UpdateStatus(result.OrderID, entity.OrderStatusReceived)
	if err != nil {
		return err
	}

	err = p.paymentRepository.UpdatePaymentStatusWithOrderID(result.OrderID, entity.PaymentStatusApproved)
	if err != nil {
		return err
	}

	return nil
}

func (p *Payments) CreatePayment(input input.CreatePaymentInput) (*output.CreatePaymentOutput, error) {
	qrcode, err := p.paymentService.GeneratePaymentQrCode(&input)
	if err != nil {
		return nil, err
	}

	payment := &entity.Payment{
		Amount:   input.Amount,
		Status:   entity.PaymentStatusPending,
		Provider: "mercadopago",
	}

	created, err := p.paymentRepository.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	err = p.orderRepository.AddPayment(input.OrderID, created.ID)
	if err != nil {
		return nil, err
	}

	return &output.CreatePaymentOutput{
		PaymentID: created.ID,
		OrderID:   qrcode.OrderID,
		QRCode:    qrcode.QRCode,
		QRCodeB64: qrcode.QRCodeB64,
	}, nil
}
