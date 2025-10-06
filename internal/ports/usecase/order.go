package usecase

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type Order interface {
	CreateOrder(input entity.Order) (int, error)
	UpdateOrderStatus(id int, status entity.OrderStatus) error
	GetOrderByID(id int) (entity.Order, error)
	ListOrders(filter input.OrderFilterInput) ([]entity.Order, error)
}
