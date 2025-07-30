package variants

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductDetails struct {
	Code     string          `json:"code"`
	Price    float64         `json:"price"`
	Category string          `json:"category"`
	Variants []VariantDetail `json:"variants"`
}

type VariantDetail struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductReader interface {
	GetProductByID(id uint) (*models.Product, error)
}

type VariantHandler struct {
	repo ProductReader
}

func NewVariantHandler(r ProductReader) *VariantHandler {
	return &VariantHandler{repo: r}
}

func (h *VariantHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/catalog/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetProductByID(uint(id))
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	details := ProductDetails{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: "",
	}

	if product.Category != nil {
		details.Category = product.Category.Name
	}

	for _, v := range product.Variants {
		price := v.Price
		if price.IsZero() {
			price = product.Price
		}
		details.Variants = append(details.Variants, VariantDetail{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}
