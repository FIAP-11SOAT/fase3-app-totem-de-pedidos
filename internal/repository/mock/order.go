package mock

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

type OrderMock struct {
	CreateOrderFunc  func(entity.Order) (int, error)
	UpdateStatusFunc func(int, entity.OrderStatus) error
	GetOrderByIDFunc func(int) (entity.Order, error)
	ListOrdersFunc   func(input.OrderFilterInput) ([]entity.Order, error)
	AddPaymentFunc   func(int, int) error
}

func NewOrderRepositoryMock() *OrderMock {
	return &OrderMock{
		CreateOrderFunc:  func(entity.Order) (int, error) { return 0, nil },
		UpdateStatusFunc: func(int, entity.OrderStatus) error { return nil },
		GetOrderByIDFunc: func(int) (entity.Order, error) { return entity.Order{}, nil },
		ListOrdersFunc:   func(input.OrderFilterInput) ([]entity.Order, error) { return []entity.Order{}, nil },
		AddPaymentFunc:   func(int, int) error { return nil },
	}
}

func (o *OrderMock) CreateOrder(order entity.Order) (int, error) { return o.CreateOrderFunc(order) }

func (o *OrderMock) UpdateStatus(i int, e entity.OrderStatus) error { return o.UpdateStatusFunc(i, e) }

func (o *OrderMock) GetOrderByID(i int) (entity.Order, error) { return o.GetOrderByIDFunc(i) }

func (o *OrderMock) ListOrders(ofi input.OrderFilterInput) ([]entity.Order, error) {
	return o.ListOrdersFunc(ofi)
}

func (o *OrderMock) AddPayment(i int, ii int) error { return o.AddPaymentFunc(i, ii) }
