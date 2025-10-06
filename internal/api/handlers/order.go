package handlers

import (
	"net/http"
	"strconv"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/dto"

	"github.com/labstack/echo/v4"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/mapper"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
)

type OrderHanlder struct {
	orderService usecase.Order
}

func NewOrderHandler(service usecase.Order) *OrderHanlder {
	return &OrderHanlder{
		orderService: service,
	}
}

func (o *OrderHanlder) CreateOrder(c echo.Context) error {
	var orderInput input.OrderInput
	if err := c.Bind(&orderInput); err != nil {
		c.Logger().Error("failed to create order", err)
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: "invalid input"})
	}

	if err := orderInput.Validate(); err != nil {
		c.Logger().Error("failed to create order", err)
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: err.Error()})
	}

	orderEntity := mapper.ToOrderDomain(orderInput)

	orderId, err := o.orderService.CreateOrder(orderEntity)
	if err != nil {
		c.Logger().Error("failed to create order", err)
		return c.JSON(http.StatusInternalServerError, dto.HttpResponseError{Error: "failed to create order"})
	}

	return c.JSON(
		http.StatusCreated,
		map[string]string{"message": "order created successfully", "orderId": strconv.Itoa(orderId)},
	)
}

func (o *OrderHanlder) UpdateOrder(c echo.Context) error {
	var content input.UpdateOrderInput
	if err := c.Bind(&content); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: "invalid input"})
	}

	if err := content.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: err.Error()})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: "invalid id"})
	}

	if err := o.orderService.UpdateOrderStatus(id, content.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HttpResponseError{Error: "could not update order"})
	}

	return c.JSON(http.StatusOK, dto.HttpResponse{Message: "order updated successfully"})
}

func (o *OrderHanlder) GetOrderById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: "invalid order ID"})
	}

	order, err := o.orderService.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.HttpResponseError{Error: "order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (o *OrderHanlder) GetOrders(c echo.Context) error {
	var filter input.OrderFilterInput
	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HttpResponseError{Error: "invalid query parameters"})
	}

	orders, err := o.orderService.ListOrders(filter)
	if err != nil {
		c.Logger().Error("failed to fetch orders", err)
		return c.JSON(http.StatusInternalServerError, dto.HttpResponseError{Error: "could not fetch orders"})
	}

	return c.JSON(http.StatusOK, orders)
}
