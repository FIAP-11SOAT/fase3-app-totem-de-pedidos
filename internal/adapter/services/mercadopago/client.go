package mercadopago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago/request"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/services/mercadopago/response"
)

var client = &http.Client{}

type MPStoreClient struct {
	AccessToken   string
	UserID        string
	ExternalPosID string
}

func NewMPStoreClient(AccessToken, UserID, ExternalPosID string) *MPStoreClient {
	return &MPStoreClient{
		AccessToken:   AccessToken,
		UserID:        UserID,
		ExternalPosID: ExternalPosID,
	}
}

func (s *MPStoreClient) GetPaymentByID(paymentID string) (*response.PaymentIDResponse, error) {
	endpoint := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responsePayment response.PaymentIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayment); err != nil {
		return nil, err
	}

	return &responsePayment, nil
}

func (s *MPStoreClient) CreateDynamicQR(input *request.CreateDynamicQRInput) (*response.DynamicQRResponse, error) {
	endpoint := fmt.Sprintf(
		"https://api.mercadopago.com/instore/orders/qr/seller/collectors/%s/pos/%s/qrs",
		s.UserID,
		s.ExternalPosID,
	)

	content, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseQrCode response.DynamicQRResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseQrCode); err != nil {
		return nil, err
	}

	return &responseQrCode, nil
}
