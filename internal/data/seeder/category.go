package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
)

func CategorySeeds() []entity.Category {
	return []entity.Category{
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
}
