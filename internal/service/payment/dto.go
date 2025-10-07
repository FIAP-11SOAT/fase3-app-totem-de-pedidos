package payment

type PaymentInput struct {
	OrderID string
	Amount  float64
	Title   string
}

type PaymentOutput struct {
	OrderID   string
	QRCode    string
	QRCodeB64 string
}
