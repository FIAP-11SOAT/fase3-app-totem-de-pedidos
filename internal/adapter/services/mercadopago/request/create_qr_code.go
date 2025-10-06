package request

type Items struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	UnitPrice   float64 `json:"unit_price"`
	TotalAmount float64 `json:"total_amount"`
}

type CreateDynamicQRInput struct {
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	ExternalReference string  `json:"external_reference"`
	NotificationURL   string  `json:"notification_url"`
	TotalAmount       float64 `json:"total_amount"`
	Items             []Items `json:"items"`
}
