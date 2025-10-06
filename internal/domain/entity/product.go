package entity

import "time"

type ProductCategory struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	ID              int
	Name            string
	Description     string
	Price           float64
	ImageURL        string
	PreparationTime int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CategoryID      int
}
