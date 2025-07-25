package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Total    int       `json:"total"`
	Products []Product `json:"products"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo models.ProductReader
}

func NewCatalogHandler(r models.ProductReader) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Apply defaults and bounds
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Filters
	category := r.URL.Query().Get("category")
	priceLt := 0.0
	if p := r.URL.Query().Get("price_lt"); p != "" {
		if parsed, err := strconv.ParseFloat(p, 64); err == nil {
			priceLt = parsed
		}
	}

	// Fetch products
	res, total, err := h.repo.GetAllProducts(offset, limit, category, priceLt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		categoryName := ""
		if p.Category != nil {
			categoryName = p.Category.Name
		}

		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: categoryName,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Total:    int(total),
		Products: products,
	})
}
