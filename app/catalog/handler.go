package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
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
	res, err := h.repo.GetAllProducts()
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
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: categoryName,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Response{Products: products}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
