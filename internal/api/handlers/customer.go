package handlers

import (
	"net/http"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	CustomerService usecase.Customer
}

func NewCustomerHandler(service usecase.Customer) *CustomerHandler {
	return &CustomerHandler{
		CustomerService: service,
	}
}

func (h *CustomerHandler) CreateCustomer(c echo.Context) error {
	var customerCreateRequest usecase.CustomerInput
	if err := c.Bind(&customerCreateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "error binding create customer request"})
	}

	customerResponse, err := h.CustomerService.CreateCustomer(&customerCreateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, customerResponse)
}

func (h *CustomerHandler) IdentifyCustomer(c echo.Context) error {
	taxId := c.QueryParam("taxid")

	if taxId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "param taxid missing"})
	}

	customerResponse, err := h.CustomerService.IdentifyCustomer(&taxId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if customerResponse == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "customer not found"})
	}

	return c.JSON(http.StatusOK, customerResponse)
}
