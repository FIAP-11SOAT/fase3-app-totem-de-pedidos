package repositories

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type Order interface {
	CreateOrder(entity.Order) (int, error)
	UpdateStatus(orderID int, status entity.OrderStatus) error
	AddPayment(orderID int, paymentID int) error
	GetOrderByID(orderID int) (entity.Order, error)
	ListOrders(filter input.OrderFilterInput) ([]entity.Order, error)
}
