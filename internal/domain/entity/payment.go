package entity

import "time"

type PaymentStatus string

type Payment struct {
	ID          int
	Amount      float64
	PaymentDate time.Time
	Status      PaymentStatus
	Provider    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusRejected PaymentStatus = "rejected"
)
