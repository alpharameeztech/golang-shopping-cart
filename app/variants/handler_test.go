package variants

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Mock repository implementing ProductReader
type mockProductRepo struct {
	product *models.Product
	err     error
}

func (m *mockProductRepo) GetProductByID(id uint) (*models.Product, error) {
	return m.product, m.err
}

func TestHandleGetByID_Success(t *testing.T) {
	product := &models.Product{
		ID:    1,
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: &models.Category{
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(12.99)},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.Zero}, // Should inherit
		},
	}

	repo := &mockProductRepo{product: product}
	handler := NewVariantHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/catalog/1", nil)
	w := httptest.NewRecorder()
	handler.HandleGetByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var body ProductDetails
	err := json.NewDecoder(res.Body).Decode(&body)
	assert.NoError(t, err)

	assert.Equal(t, "PROD001", body.Code)
	assert.Equal(t, 10.99, body.Price)
	assert.Equal(t, "Clothing", body.Category)
	assert.Len(t, body.Variants, 2)
	assert.Equal(t, 12.99, body.Variants[0].Price)
	assert.Equal(t, 10.99, body.Variants[1].Price) // inherited
}

func TestHandleGetByID_NotFound(t *testing.T) {
	repo := &mockProductRepo{product: nil, err: errors.New("not found")}
	handler := NewVariantHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/catalog/999", nil)
	w := httptest.NewRecorder()
	handler.HandleGetByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestHandleGetByID_InvalidID(t *testing.T) {
	repo := &mockProductRepo{}
	handler := NewVariantHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/catalog/abc", nil)
	w := httptest.NewRecorder()
	handler.HandleGetByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestHandleGetByID_ZeroID(t *testing.T) {
	repo := &mockProductRepo{}
	handler := NewVariantHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/catalog/0", nil)
	w := httptest.NewRecorder()
	handler.HandleGetByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestHandleGetByID_InvalidPath(t *testing.T) {
	repo := &mockProductRepo{}
	handler := NewVariantHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/wrongpath/1", nil)
	w := httptest.NewRecorder()
	handler.HandleGetByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Since we only trim prefix in actual handler, not checking for segment count here
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
