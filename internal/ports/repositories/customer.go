package repositories

import "github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"

type Customer interface {
	CreateCustomer(customer *entity.Customer) (*entity.Customer, error)
	IdentifyCustomer(taxID *string) (*entity.Customer, error)
}
