package input

type CategoryInput struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CategoryInput) Validate() error {
	return nil
}
