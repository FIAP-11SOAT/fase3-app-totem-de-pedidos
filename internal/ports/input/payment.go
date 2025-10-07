package input

type CreatePaymentInput struct {
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Title   string  `json:"title"`
}
