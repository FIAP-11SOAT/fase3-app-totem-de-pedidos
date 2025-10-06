package usecase

import (
	"errors"
	"testing"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) CreateCustomer(customer *entity.Customer) (*entity.Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) IdentifyCustomer(taxID *string) (*entity.Customer, error) {
	args := m.Called(taxID)
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func TestCreateCustomer_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	useCase := NewCustomerUseCase(mockRepo)

	input := &usecase.CustomerInput{
		Nome:  "João",
		Email: "joao@example.com",
		TaxID: "12345678909",
	}

	expected := &entity.Customer{
		Name:  "João",
		Email: "joao@example.com",
		TaxID: "12345678909",
	}

	mockRepo.On("CreateCustomer", mock.AnythingOfType("*entity.Customer")).Return(expected, nil)

	result, err := useCase.CreateCustomer(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer_InvalidCPF(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	useCase := NewCustomerUseCase(mockRepo)

	input := &usecase.CustomerInput{
		Nome:  "Maria",
		Email: "maria@example.com",
		TaxID: "123",
	}

	result, err := useCase.CreateCustomer(input)

	assert.Nil(t, result)
	assert.EqualError(t, err, "CPF inválido")
}

func TestIdentifyCustomer_Success(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	useCase := NewCustomerUseCase(mockRepo)

	cpf := "12345678909"
	expected := &entity.Customer{
		Name:  "Ana",
		Email: "ana@example.com",
		TaxID: cpf,
	}

	mockRepo.On("IdentifyCustomer", &cpf).Return(expected, nil)

	result, err := useCase.IdentifyCustomer(&cpf)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestIdentifyCustomer_InvalidCPF(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	useCase := NewCustomerUseCase(mockRepo)

	cpf := "000"
	result, err := useCase.IdentifyCustomer(&cpf)

	assert.Nil(t, result)
	assert.EqualError(t, err, "CPF inválido")
}

func TestIdentifyCustomer_ErrorFromRepo(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	useCase := NewCustomerUseCase(mockRepo)

	cpf := "12345678909"
	mockRepo.On("IdentifyCustomer", &cpf).Return(&entity.Customer{}, errors.New("db error"))

	result, err := useCase.IdentifyCustomer(&cpf)

	assert.Nil(t, result)
	assert.EqualError(t, err, "error getting customer")
}
