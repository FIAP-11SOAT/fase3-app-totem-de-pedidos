package entity

import "time"

type OrderStatus string

const (
	OrderStatusPaymentPending OrderStatus = "payment_pending"
	OrderStatusReceived       OrderStatus = "received"
	OrderStatusPreparing      OrderStatus = "preparing"
	OrderStatusReady          OrderStatus = "ready"
	OrderStatusCompleted      OrderStatus = "completed"
)

type Order struct {
	ID                   int
	OrderDate            time.Time
	NotificationAttempts int
	Status               OrderStatus
	TotalAmount          float64
	CreatedAt            time.Time
	CustomerID           *int
	Items                []OrderItem
}

type OrderItem struct {
	ID        int
	Quantity  int
	Price     float64
	CreatedAt time.Time
	OrderID   int
	ProductID int
}
