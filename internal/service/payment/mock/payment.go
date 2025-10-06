package mock

import "github.com/FIAP-11SOAT/totem-de-pedidos/internal/service/payment"

type PaymentServiceMock struct {
	PaymentFunc func(payment.PaymentInput) (payment.PaymentOutput, error)
}

func NewPaymentServiceMock() *PaymentServiceMock {
	return &PaymentServiceMock{
		PaymentFunc: func(payment.PaymentInput) (payment.PaymentOutput, error) { return payment.PaymentOutput{}, nil },
	}
}

func (p *PaymentServiceMock) Payment(input payment.PaymentInput) (payment.PaymentOutput, error) {
	return p.PaymentFunc(input)
}
