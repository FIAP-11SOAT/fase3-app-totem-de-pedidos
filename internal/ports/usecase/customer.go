package usecase

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
)

type CustomerInput struct {
	ID    int    `json:"id,omitempty"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
	TaxID string `json:"taxid"`
}

type Customer interface {
	CreateCustomer(customerDTO *CustomerInput) (*entity.Customer, error)
	IdentifyCustomer(taxID *string) (*entity.Customer, error)
}
