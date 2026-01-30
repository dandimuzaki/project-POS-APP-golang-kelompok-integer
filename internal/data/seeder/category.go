package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func CategorySeeds(db *gorm.DB) ([]entity.Category, error) {
	var count int64
	db.Model(&entity.Category{}).Count(&count)

	var categories []entity.Category
	if count > 0 {
		db.Find(&categories)
		return categories, nil
	}

	categories = []entity.Category{
		{
			Name:        "Coffee",
			Description: "Berbagai macam kopi panas dan dingin",
		},
		{
			Name:        "Non Coffee",
			Description: "Minuman non kopi seperti teh dan coklat",
		},
		{
			Name:        "Food",
			Description: "Makanan berat dan ringan",
		},
		{
			Name:        "Snack",
			Description: "Camilan pendamping minum kopi",
		},
		{
			Name:        "Dessert",
			Description: "Makanan penutup",
		},
	}

	if err := db.Create(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
