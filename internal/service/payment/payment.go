package payment

import (
	"os"
	"strconv"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago/request"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/qrcode"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/output"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/service"
)

type MPPaymentService struct {
	client *mercadopago.MPStoreClient
}

func NewMPService() service.PaymentService {
	client := mercadopago.NewMPStoreClient(
		os.Getenv("MP_TOKEN"),
		os.Getenv("MP_USER_ID"),
		os.Getenv("MP_EXTERNAL_POS_ID"),
	)
	return &MPPaymentService{
		client: client,
	}
}

func (p *MPPaymentService) GetPaymentByID(paymentID string) (*output.GetPaymentOutput, error) {
	payment, err := p.client.GetPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	orderID, err := strconv.Atoi(payment.ExternalReference)
	if err != nil {
		return nil, err
	}

	return &output.GetPaymentOutput{
		PaymentID: paymentID,
		OrderID:   orderID,
		Status:    payment.Status,
	}, nil
}

func (p *MPPaymentService) GeneratePaymentQrCode(input *input.CreatePaymentInput) (*output.CreatePaymentOutput, error) {
	payload := &request.CreateDynamicQRInput{
		Title:             input.Title,
		Description:       input.Title,
		ExternalReference: strconv.Itoa(input.OrderID),
		NotificationURL:   os.Getenv("MP_NOTIFICATION_URL"),
		TotalAmount:       input.Amount,
		Items: []request.Items{
			{
				ID:          "Totem-Payment",
				Title:       input.Title,
				Description: input.Title,
				Category:    "Totem-Payment",
				Quantity:    1,
				UnitMeasure: "unit",
				UnitPrice:   input.Amount,
				TotalAmount: input.Amount,
			},
		},
	}

	result, err := p.client.CreateDynamicQR(payload)
	if err != nil {
		return nil, err
	}

	qrcodeB64, err := qrcode.CreateQRCode(result.QRData)
	if err != nil {
		return nil, err
	}

	return &output.CreatePaymentOutput{
		QRCodeB64: qrcodeB64,
		QRCode:    result.QRData,
		OrderID:   input.OrderID,
	}, nil
}
