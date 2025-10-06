package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/repositories"
)

type order struct {
	orderRepo   repositories.Order
	productRepo repositories.Product
}

func NewOrderUseCase(
	orderRepo repositories.Order,
	productRepo repositories.Product,
) *order {
	return &order{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (o *order) CreateOrder(input entity.Order) (int, error) {
	if len(input.Items) == 0 {
		return -1, errors.New("order must have at least one item")
	}

	var total float64
	for i, item := range input.Items {
		product, err := o.productRepo.FindProductByID(context.Background(), strconv.Itoa(item.ProductID))
		if err != nil {
			return -1, err
		}

		item.Price = product.Price * float64(item.Quantity)
		total += item.Price

		input.Items[i] = item
	}

	input.TotalAmount = total
	input.Status = entity.OrderStatusPaymentPending

	orderId, err := o.orderRepo.CreateOrder(input)
	if err != nil {
		return -1, err
	}

	return orderId, err
}

func (o *order) UpdateOrderStatus(id int, status entity.OrderStatus) error {
	return o.orderRepo.UpdateStatus(id, status)
}

func (o *order) GetOrderByID(id int) (entity.Order, error) {
	return o.orderRepo.GetOrderByID(id)
}

func (o *order) ListOrders(filter input.OrderFilterInput) ([]entity.Order, error) {
	if filter.CustomerID != nil {
		orders, err := o.orderRepo.ListOrders(filter)
		if err != nil {
			return []entity.Order{}, nil
		}

		for i := range orders {
			orders[i].TotalAmount = orders[i].TotalAmount * 0.95 // aplica 5% de desconto para o cliente
		}

		return orders, nil
	}

	return o.orderRepo.ListOrders(filter)
}
