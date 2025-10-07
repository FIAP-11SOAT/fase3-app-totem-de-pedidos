package output

import "time"

type ProductOutput struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	ImageURL        string    `json:"image_url"`
	PreparationTime int       `json:"preparation_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CategoryID      int       `json:"category_id"`
}
