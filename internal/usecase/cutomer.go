package usecase

import (
	"fmt"

	"github.com/klassmann/cpfcnpj"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
)

type Customer struct {
	CustomerRepository repositories.Customer
}

func NewCustomerUseCase(repository repositories.Customer) usecase.Customer {
	return &Customer{CustomerRepository: repository}
}

func (p *Customer) CreateCustomer(customerDTO *usecase.CustomerInput) (*entity.Customer, error) {
	customerEntity := entity.Customer{
		Name:  customerDTO.Nome,
		Email: customerDTO.Email,
		TaxID: customerDTO.TaxID,
	}

	if !cpfcnpj.ValidateCPF(customerDTO.TaxID) {
		return nil, fmt.Errorf("CPF inválido")
	}

	customerResponse, err := p.CustomerRepository.CreateCustomer(&customerEntity)
	if err != nil {
		return nil, err
	}
	return customerResponse, nil
}

func (p *Customer) IdentifyCustomer(taxID *string) (*entity.Customer, error) {
	if !cpfcnpj.ValidateCPF(*taxID) {
		return nil, fmt.Errorf("CPF inválido")
	}

	customerResponse, err := p.CustomerRepository.IdentifyCustomer(taxID)
	if err != nil {
		return nil, fmt.Errorf("error getting customer")
	}
	return customerResponse, err
}
