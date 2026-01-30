package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func ProductSeeds(db *gorm.DB, c []entity.Category) ([]entity.Product, error) {
	if len(c) < 2 {
		return nil, fmt.Errorf("not enough categories to seed products")
	}

	var products []entity.Product
	var count int64
	db.Model(&entity.Product{}).Count(&count)
	if count > 0 {
		db.Find(&products)
		return products, nil
	}
	
	products = []entity.Product{
		{
			Name:        "Espresso",
			Description: "Kopi hitam pekat tanpa gula",
			Price:       25000,
			Stock:       100,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/espresso.jpg",
			CategoryID:  c[0].ID, // Coffee
		},
		{
			Name:        "Cappuccino",
			Description: "Kopi dengan susu dan foam",
			Price:       30000,
			Stock:       80,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/cappuccino.jpg",
			CategoryID:  c[0].ID, // Coffee
		},
		{
			Name:        "Matcha Latte",
			Description: "Minuman matcha dengan susu",
			Price:       28000,
			Stock:       60,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/matcha.jpg",
			CategoryID:  c[1].ID, // Non Coffee
		},
	}

	if err := db.Create(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
