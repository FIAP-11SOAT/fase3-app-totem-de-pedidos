package entity

import "time"

type Customer struct {
	ID        int
	Name      string
	Email     string
	TaxID     string
	CreatedAt time.Time
}
