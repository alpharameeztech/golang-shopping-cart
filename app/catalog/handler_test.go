package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	_ "strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	products []models.Product
	total    int64
	err      error
}

func (m *mockRepo) GetProductByID(id uint) (*models.Product, error) {
	return nil, nil
}

func (m *mockRepo) GetAllProducts(offset, limit int, category string, priceLt float64) ([]models.Product, int64, error) {
	return m.products, m.total, m.err
}

func TestHandleGet_Success(t *testing.T) {
	h := NewCatalogHandler(&mockRepo{
		products: []models.Product{
			{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(9.99),
				Category: &models.Category{
					Name: "Shoes",
				},
			},
		},
		total: 1,
	})

	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=0&limit=5", nil)
	rr := httptest.NewRecorder()
	h.HandleGet(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var resp Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Len(t, resp.Products, 1)
	assert.Equal(t, "PROD001", resp.Products[0].Code)
	assert.Equal(t, 9.99, resp.Products[0].Price)
	assert.Equal(t, "Shoes", resp.Products[0].Category)
}

func TestHandleGet_DefaultsAndBounds(t *testing.T) {
	h := NewCatalogHandler(&mockRepo{})
	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=-5&limit=200", nil)
	rr := httptest.NewRecorder()
	h.HandleGet(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandleGet_WithFilters(t *testing.T) {
	h := NewCatalogHandler(&mockRepo{
		products: []models.Product{
			{
				Code:  "PROD002",
				Price: decimal.NewFromFloat(4.5),
				Category: &models.Category{
					Name: "Accessories",
				},
			},
		},
		total: 1,
	})

	req := httptest.NewRequest(http.MethodGet, "/catalog?category=Accessories&price_lt=5", nil)
	rr := httptest.NewRecorder()
	h.HandleGet(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "Accessories", resp.Products[0].Category)
	assert.Equal(t, 4.5, resp.Products[0].Price)
}

func TestHandleGet_InternalError(t *testing.T) {
	h := NewCatalogHandler(&mockRepo{
		err: fmt.Errorf("database error"),
	})

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rr := httptest.NewRecorder()
	h.HandleGet(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}
