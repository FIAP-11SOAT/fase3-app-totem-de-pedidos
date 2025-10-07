package mapper

import (
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
)

func ToOrderDomain(input input.OrderInput, customerId string) entity.Order {
	items := make([]entity.OrderItem, len(input.Items))

	for i, item := range input.Items {
		items[i] = entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return entity.Order{
		CustomerID:           &customerId,
		Items:                items,
		Status:               "PENDING",
		NotificationAttempts: 0,
	}
}
