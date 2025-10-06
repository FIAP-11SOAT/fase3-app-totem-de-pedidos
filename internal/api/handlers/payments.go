package handlers

import (
	"fmt"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"

	"github.com/labstack/echo/v4"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/dto"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/mapper"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
)

type PaymentHandler struct {
	paymentService usecase.Payments
}

func NewPaymentHandler(paymentService usecase.Payments) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) GetByID(c echo.Context) error {
	paymentID := c.Param("id")
	if paymentID == "" {
		return c.JSON(400, dto.HttpResponse{Message: "Payment ID is required"})
	}

	payment, err := h.paymentService.GetPaymentByID(paymentID)
	if err != nil {
		c.Logger().Error(dto.LoggerInfo{Scope: "get:payment:id", Message: "Failed to get payment by ID"})
		return c.JSON(500, dto.HttpResponse{Message: "Internal server error"})
	}

	return c.JSON(200, mapper.MapPaymentToOutput(payment))
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var data input.CreatePaymentInput
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(dto.LoggerInfo{Scope: "post:payment", Message: "Failed to bind payment input"})
		return c.JSON(400, dto.HttpResponse{Message: "Invalid input"})
	}

	payment, err := h.paymentService.CreatePayment(data)
	if err != nil {
		c.Logger().Error(dto.LoggerInfo{Scope: "post:payment", Message: "Failed to create payment"})
		return c.JSON(500, dto.HttpResponse{Message: "Internal server error"})
	}

	return c.JSON(200, payment)
}

func (h *PaymentHandler) PaymentWebHook(c echo.Context) error {

	var queryType = c.QueryParam("type")
	if queryType != "payment" {
		return c.JSON(202, dto.HttpResponse{Message: "Invalid Params Type"})
	}

	// Ambiente de desenvolvimento nao fornece a assinatura correta somente na simulaçao do webhook na pagina
	//isValid := mercadopago.CheckWebhookSignature(mercadopago.CheckWebhookSignatureInput{
	//	XSignature: c.Request().Header.Get("X-Signature"),
	//	XRequestId: c.Request().Header.Get("X-Request-Id"),
	//	DataID:     c.QueryParam("data.id"),
	//	Secret:     os.Getenv("MP_WEBHOOK_SECRET"),
	//})
	//if !isValid {
	//	c.Logger().Error(dto.LoggerInfo{Scope: "post:payment:webhook", Message: "Invalid webhook signature"})
	//	return c.NoContent(202)
	//}

	var payload mercadopago.WebhookPayload
	if err := c.Bind(&payload); err != nil {
		fmt.Println(err)
		c.Logger().Error(dto.LoggerInfo{Scope: "post:payment:webhook", Message: "Failed to bind webhook payload"})
		return c.NoContent(202)
	}

	err := h.paymentService.PaymentWebHook(&payload)
	if err != nil {
		c.Logger().Error(dto.LoggerInfo{Scope: "post:payment:webhook", Message: "Failed to process webhook payload" + err.Error()})
		return c.NoContent(202)
	}

	return c.JSON(200, dto.HttpResponse{
		Message: "Webhook processed successfully",
	})
}
