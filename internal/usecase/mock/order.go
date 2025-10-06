package mock

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type orderServiceMock struct {
	CreateOrderFunc       func(entity.Order) (int, error)
	UpdateOrderStatusFunc func(int, entity.OrderStatus) error
	GetOrderByIDFunc      func(int) (entity.Order, error)
	ListOrdersFunc        func(input.OrderFilterInput) ([]entity.Order, error)
	CheckoutFunc          func(int) error
}

func NewOrderServiceMock() *orderServiceMock {
	return &orderServiceMock{
		CreateOrderFunc:       func(o entity.Order) (int, error) { return -1, nil },
		UpdateOrderStatusFunc: func(int, entity.OrderStatus) error { return nil },
		GetOrderByIDFunc:      func(i int) (entity.Order, error) { return entity.Order{}, nil },
		ListOrdersFunc:        func(ofi input.OrderFilterInput) ([]entity.Order, error) { return []entity.Order{}, nil },
		CheckoutFunc:          func(i int) error { return nil },
	}
}

func (m *orderServiceMock) CreateOrder(o entity.Order) (int, error) {
	return m.CreateOrderFunc(o)
}

func (m *orderServiceMock) UpdateOrderStatus(i int, e entity.OrderStatus) error {
	return m.UpdateOrderStatusFunc(i, e)
}

func (m *orderServiceMock) GetOrderByID(i int) (entity.Order, error) { return m.GetOrderByIDFunc(i) }

func (m *orderServiceMock) ListOrders(filter input.OrderFilterInput) ([]entity.Order, error) {
	return m.ListOrdersFunc(filter)
}

func (m *orderServiceMock) Checkout(i int) error { return m.CheckoutFunc(i) }
