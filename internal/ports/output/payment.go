package output

import "time"

type PaymentOutput struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Status      string    `json:"status"`
	Provider    string    `json:"provider"`
}

type CreatePaymentOutput struct {
	PaymentID int    `json:"payment_id"`
	OrderID   int    `json:"order_id"`
	QRCode    string `json:"qrcode"`
	QRCodeB64 string `json:"qrcode_b64"`
}

type GetPaymentOutput struct {
	PaymentID string `json:"payment_id"`
	OrderID   int    `json:"order_id"`
	Status    string `json:"status"`
}
