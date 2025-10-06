package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api/handlers"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) CreateCustomer(input *usecase.CustomerInput) (*entity.Customer, error) {
	args := m.Called(input)
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerService) IdentifyCustomer(taxId *string) (*entity.Customer, error) {
	args := m.Called(taxId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func TestCreateCustomer_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"João","email":"joao@email.com","tax_id":"12345678901"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockCustomerService)
	handler := handlers.NewCustomerHandler(mockService)

	mockService.On("CreateCustomer", mock.Anything).Return(&entity.Customer{Name: "João"}, nil)

	err := handler.CreateCustomer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCreateCustomer_BindError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockCustomerService)
	handler := handlers.NewCustomerHandler(mockService)

	err := handler.CreateCustomer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIdentifyCustomer_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?taxid=123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockCustomerService)
	handler := handlers.NewCustomerHandler(mockService)

	taxid := "123"
	mockService.On("IdentifyCustomer", &taxid).Return(&entity.Customer{TaxID: "123"}, nil)

	err := handler.IdentifyCustomer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestIdentifyCustomer_NoTaxID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockCustomerService)
	handler := handlers.NewCustomerHandler(mockService)

	err := handler.IdentifyCustomer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIdentifyCustomer_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?taxid=123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockCustomerService)
	handler := handlers.NewCustomerHandler(mockService)

	taxid := "123"
	mockService.On("IdentifyCustomer", &taxid).Return(nil, nil)

	err := handler.IdentifyCustomer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
