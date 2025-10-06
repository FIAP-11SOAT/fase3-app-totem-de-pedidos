package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetCategories(ctx context.Context) ([]*entity.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Category), args.Error(1)
}
func (m *MockCategoryRepository) CreateCategory(ctx context.Context, cat *entity.Category) (int, error) {
	args := m.Called(ctx, cat)
	return args.Int(0), args.Error(1)
}
func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, cat *entity.Category) (*entity.Category, error) {
	args := m.Called(ctx, cat)
	return args.Get(0).(*entity.Category), args.Error(1)
}
func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockCategoryRepository) FindCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Category), args.Error(1)
}

func TestGetCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	expected := []*entity.Category{{ID: 1, Name: "Bebidas"}}
	mockRepo.On("GetCategories", mock.Anything).Return(expected, nil)

	uc := NewCategoryCase(mockRepo)
	cats, err := uc.GetCategories()
	assert.NoError(t, err)
	assert.Equal(t, expected, cats)
}

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	inputCat := &input.CategoryInput{Name: "Lanches", Description: "Sanduíches"}
	mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(2, nil)

	uc := NewCategoryCase(mockRepo)
	cat, err := uc.CreateCategory(inputCat)
	assert.NoError(t, err)
	assert.Equal(t, 2, cat.ID)
	assert.Equal(t, "Lanches", cat.Name)
}

func TestCreateCategory_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	inputCat := &input.CategoryInput{Name: "Lanches", Description: "Sanduíches"}
	mockRepo.On("CreateCategory", mock.Anything, mock.Anything).Return(0, errors.New("insert error"))

	uc := NewCategoryCase(mockRepo)
	cat, err := uc.CreateCategory(inputCat)
	assert.Error(t, err)
	assert.Nil(t, cat)
}

func TestUpdateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	now := time.Now()
	inputCat := &entity.Category{ID: 1, Name: "Doces", Description: "Sobremesas", CreatedAt: now}
	mockRepo.On("UpdateCategory", mock.Anything, mock.Anything).Return(inputCat, nil)

	uc := NewCategoryCase(mockRepo)
	cat, err := uc.UpdateCategory(inputCat)
	assert.NoError(t, err)
	assert.Equal(t, inputCat, cat)
}

func TestDeleteCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockRepo.On("DeleteCategory", mock.Anything, 1).Return(nil)

	uc := NewCategoryCase(mockRepo)
	err := uc.DeleteCategory(1)
	assert.NoError(t, err)
}

func TestDeleteCategory_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	mockRepo.On("DeleteCategory", mock.Anything, 1).Return(errors.New("delete error"))

	uc := NewCategoryCase(mockRepo)
	err := uc.DeleteCategory(1)
	assert.Error(t, err)
}

func TestFindCategoryByID_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	expected := &entity.Category{ID: 1, Name: "Bebidas"}
	mockRepo.On("FindCategoryByID", mock.Anything, 1).Return(expected, nil)

	uc := NewCategoryCase(mockRepo)
	cat, err := uc.FindCategoryByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, cat)
}
