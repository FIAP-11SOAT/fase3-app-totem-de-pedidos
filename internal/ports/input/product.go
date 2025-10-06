package input

type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	CategoryID  int     `json:"category_id"`
}

type ProductFilterInput struct {
	Name         string `query:"name"`
	CategoryName string `query:"category_name"`
}

func (p *ProductInput) Validate() error {
	return nil
}
