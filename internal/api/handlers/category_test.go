//go:build unit

package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/domain/entity"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/ports/input"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService implements usecase.Category for testing
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) GetCategories() ([]*entity.Category, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Category), args.Error(1)
}
func (m *MockCategoryService) FindCategoryByID(id int) (*entity.Category, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Category), args.Error(1)
}
func (m *MockCategoryService) CreateCategory(input *input.CategoryInput) (*entity.Category, error) {
	args := m.Called(input)
	return args.Get(0).(*entity.Category), args.Error(1)
}
func (m *MockCategoryService) UpdateCategory(cat *entity.Category) (*entity.Category, error) {
	args := m.Called(cat)
	return args.Get(0).(*entity.Category), args.Error(1)
}
func (m *MockCategoryService) DeleteCategory(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupEcho() (*echo.Echo, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e, rec
}

func TestListAllCategories(t *testing.T) {
	e, rec := setupEcho()
	mockService := new(MockCategoryService)
	handler := &CategoryHandler{categoryService: mockService}

	t.Run("success", func(t *testing.T) {
		cats := []*entity.Category{{ID: 1, Name: "Bebidas"}}
		mockService.On("GetCategories").Return(cats, nil)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		ctx := e.NewContext(req, rec)
		err := handler.ListAllCategories(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp []*entity.Category
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Len(t, resp, 1)
		mockService.AssertExpectations(t)
	})

	t.Run("empty", func(t *testing.T) {
		rec = httptest.NewRecorder()
		mockService.On("GetCategories").Return([]*entity.Category{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		ctx := e.NewContext(req, rec)
		err := handler.ListAllCategories(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		mockService.On("GetCategories").Return(nil, errors.New("db error"))
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		ctx := e.NewContext(req, rec)
		err := handler.ListAllCategories(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestFindCategoryByID(t *testing.T) {
	e, rec := setupEcho()
	mockService := new(MockCategoryService)
	handler := &CategoryHandler{categoryService: mockService}

	t.Run("success", func(t *testing.T) {
		cat := &entity.Category{ID: 2, Name: "Lanches"}
		mockService.On("FindCategoryByID", 2).Return(cat, nil)
		req := httptest.NewRequest(http.MethodGet, "/categories/2", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("2")
		err := handler.FindCategoryByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp entity.Category
		json.NewDecoder(rec.Body).Decode(&resp)
		assert.Equal(t, 2, resp.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		rec = httptest.NewRecorder()
		mockService.On("FindCategoryByID", 3).Return((*entity.Category)(nil), nil)
		req := httptest.NewRequest(http.MethodGet, "/categories/3", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("3")
		err := handler.FindCategoryByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("abc")
		err := handler.FindCategoryByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		mockService.On("FindCategoryByID", 4).Return((*entity.Category)(nil), errors.New("db error"))
		req := httptest.NewRequest(http.MethodGet, "/categories/4", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("4")
		err := handler.FindCategoryByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

type mockCategoryInput struct {
	input.CategoryInput
	ValidateFunc func() error
}

func (m *mockCategoryInput) Validate() error {
	return m.ValidateFunc()
}

func TestCreateCategory(t *testing.T) {
	e, rec := setupEcho()
	mockService := new(MockCategoryService)
	handler := &CategoryHandler{categoryService: mockService}

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Sobremesas"}`
		catInput := &input.CategoryInput{Name: "Sobremesas"}
		cat := &entity.Category{ID: 10, Name: "Sobremesas"}
		mockService.On("CreateCategory", catInput).Return(cat, nil)
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		// Patch Validate method
		origValidate := catInput.Validate
		catInput.Validate = func() error { return nil }
		defer func() { catInput.Validate = origValidate }()
		err := handler.CreateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		err := handler.CreateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validate error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		body := `{"name":""}`
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		// Patch Validate method
		invalidInput := &input.CategoryInput{Name: ""}
		origValidate := invalidInput.Validate
		invalidInput.Validate = func() error { return errors.New("validation error") }
		defer func() { invalidInput.Validate = origValidate }()
		err := handler.CreateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		body := `{"name":"Bebidas"}`
		catInput := &input.CategoryInput{Name: "Bebidas"}
		mockService.On("CreateCategory", catInput).Return((*entity.Category)(nil), errors.New("db error"))
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		origValidate := catInput.Validate
		catInput.Validate = func() error { return nil }
		defer func() { catInput.Validate = origValidate }()
		err := handler.CreateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateCategory(t *testing.T) {
	e, rec := setupEcho()
	mockService := new(MockCategoryService)
	handler := &CategoryHandler{categoryService: mockService}

	t.Run("success", func(t *testing.T) {
		body := `{"name":"Atualizado"}`
		cat := &entity.Category{ID: 5, Name: "Atualizado"}
		mockService.On("UpdateCategory", cat).Return(cat, nil)
		req := httptest.NewRequest(http.MethodPut, "/categories/5", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("5")
		err := handler.UpdateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/categories/6", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("6")
		err := handler.UpdateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/categories/abc", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("abc")
		err := handler.UpdateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		body := `{"name":"Falha"}`
		cat := &entity.Category{ID: 7, Name: "Falha"}
		mockService.On("UpdateCategory", cat).Return((*entity.Category)(nil), errors.New("db error"))
		req := httptest.NewRequest(http.MethodPut, "/categories/7", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("7")
		err := handler.UpdateCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	e, rec := setupEcho()
	mockService := new(MockCategoryService)
	handler := &CategoryHandler{categoryService: mockService}

	t.Run("success", func(t *testing.T) {
		mockService.On("DeleteCategory", 8).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/categories/8", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("8")
		err := handler.DeleteCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		rec = httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("abc")
		err := handler.DeleteCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		rec = httptest.NewRecorder()
		mockService.On("DeleteCategory", 9).Return(errors.New("db error"))
		req := httptest.NewRequest(http.MethodDelete, "/categories/9", nil)
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("9")
		err := handler.DeleteCategory(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
