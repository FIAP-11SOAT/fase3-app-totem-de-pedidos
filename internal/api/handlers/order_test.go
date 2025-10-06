package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/usecase/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderHandler(t *testing.T) {
	t.Run("should return 201 when order is created successfully", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.CreateOrderFunc = func(order entity.Order) (int, error) { return 42, nil }

		handler := handlers.NewOrderHandler(mockService)

		var cId int = 1
		body := input.OrderInput{
			CustomerID: &cId,
			Items: []input.OrderItemInput{
				{ProductID: 1, Quantity: 2},
			},
		}

		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "order created successfully", resp["message"])
		assert.Equal(t, "42", resp["orderId"])
	})

	t.Run("should return 400 on invalid json body", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(`invalid`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 500 on service error", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.CreateOrderFunc = func(order entity.Order) (int, error) {
			return 0, errors.New("something went wrong")
		}

		handler := handlers.NewOrderHandler(mockService)

		var cId int = 1
		body := input.OrderInput{
			CustomerID: &cId,
			Items: []input.OrderItemInput{
				{ProductID: 1, Quantity: 2},
			},
		}

		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateOrderHandler(t *testing.T) {
	t.Run("should return 200 when status is updated successfully", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.UpdateOrderStatusFunc = func(i int, os entity.OrderStatus) error { return nil }

		handler := handlers.NewOrderHandler(mockService)

		body := map[string]string{"status": "COMPLETED"}
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/order/1", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err = handler.UpdateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "order updated successfully", resp["message"])
	})

	t.Run("should return 400 on invalid json body", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/order/1", bytes.NewReader([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 500 when service returns error", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.UpdateOrderStatusFunc = func(i int, os entity.OrderStatus) error { return errors.New("failed to update") }

		handler := handlers.NewOrderHandler(mockService)

		body := map[string]string{"status": "CANCELLED"}
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, "/order/999", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		err = handler.UpdateOrder(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetOrderByIdHandler(t *testing.T) {
	t.Run("should return 200 when order is found", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.GetOrderByIDFunc = func(id int) (entity.Order, error) {
			var cId int = 1
			return entity.Order{
				ID:          id,
				CustomerID:  &cId,
				Status:      "PENDING",
				TotalAmount: 50.0,
				Items: []entity.OrderItem{
					{ProductID: 1, Quantity: 2, Price: 25.0},
				},
			}, nil
		}

		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/order/1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.GetOrderById(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return 400 on invalid id", func(t *testing.T) {
		e := echo.New()

		handler := handlers.NewOrderHandler(mock.NewOrderServiceMock())

		req := httptest.NewRequest(http.MethodGet, "/order/invalid", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetOrderById(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 404 when order is not found", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.GetOrderByIDFunc = func(id int) (entity.Order, error) {
			return entity.Order{}, errors.New("not found")
		}

		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/order/999", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/order/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.GetOrderById(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestGetOrdersHandler(t *testing.T) {
	t.Run("should return 200 with filtered orders", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.ListOrdersFunc = func(input.OrderFilterInput) ([]entity.Order, error) {
			var custumer1 int = 1
			var customer2 int = 2
			return []entity.Order{
				{ID: 1, Status: "PENDING", CustomerID: &custumer1},
				{ID: 2, Status: "PENDING", CustomerID: &customer2},
			}, nil
		}

		handler := handlers.NewOrderHandler(mockService)

		q := make(url.Values)
		q.Set("status", "PENDING")

		req := httptest.NewRequest(http.MethodGet, "/orders?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetOrders(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return 400 on invalid query params", func(t *testing.T) {
		e := echo.New()
		mockService := mock.NewOrderServiceMock()

		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/orders?customer_id=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.QueryParams().Set("customer_id", "invalid")

		err := handler.GetOrders(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 500 when service fails", func(t *testing.T) {
		e := echo.New()

		mockService := mock.NewOrderServiceMock()
		mockService.ListOrdersFunc = func(input.OrderFilterInput) ([]entity.Order, error) {
			return nil, errors.New("unexpected error")
		}

		handler := handlers.NewOrderHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/orders?status=COMPLETED", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetOrders(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
