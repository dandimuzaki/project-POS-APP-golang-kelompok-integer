package data

import "project-POS-APP-golang-integer/internal/data/entity"

func ProductSeeds() []entity.Product {
	return []entity.Product{
		{
			Name:        "Espresso",
			Description: "Kopi hitam pekat tanpa gula",
			Price:       25000,
			Stock:       100,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/espresso.jpg",
			CategoryID:  1, // Coffee
		},
		{
			Name:        "Cappuccino",
			Description: "Kopi dengan susu dan foam",
			Price:       30000,
			Stock:       80,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/cappuccino.jpg",
			CategoryID:  1, // Coffee
		},
		{
			Name:        "Matcha Latte",
			Description: "Minuman matcha dengan susu",
			Price:       28000,
			Stock:       60,
			MinStock:    10,
			Status:      entity.ProductStatusActive,
			ImageURL:    "https://cdn.posapp.com/products/matcha.jpg",
			CategoryID:  2, // Non Coffee
		},
	}
}
