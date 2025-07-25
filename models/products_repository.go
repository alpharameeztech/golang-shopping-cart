package models

import (
	"gorm.io/gorm"
)

type ProductReader interface {
	GetAllProducts(offset, limit int) ([]Product, int64, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{}).Preload("Variants").Preload("Category")
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
