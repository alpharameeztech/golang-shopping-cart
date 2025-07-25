package models

import (
	"gorm.io/gorm"
)

type ProductReader interface {
	GetAllProducts(offset, limit int, category string, priceLt float64) ([]Product, int64, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int, category string, priceLt float64) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{}).Preload("Variants").Preload("Category")

	if category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(categories.name) = LOWER(?)", category)
	}

	if priceLt > 0 {
		query = query.Where("price < ?", priceLt)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
