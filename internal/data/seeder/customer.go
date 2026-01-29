package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func CustomerSeeds(db *gorm.DB) ([]entity.Customer, error) {
	var count int64
	db.Model(&entity.Customer{}).Count(&count)

	var customers []entity.Customer
	if count > 0 {
		db.Find(&customers)
		return customers, nil
	}

	customers = []entity.Customer{
		{
			Title:     entity.CustomerTitleMr,
			FirstName: "Budi",
			LastName:  "Santoso",
			Phone:     "+628111111111",
			Email:     "budi.santoso@gmail.com",
		},
		{
			Title:     entity.CustomerTitleMs,
			FirstName: "Siti",
			LastName:  "Aisyah",
			Phone:     "+628122222222",
			Email:     "siti.aisyah@gmail.com",
		},
		{
			Title:     entity.CustomerTitleDr,
			FirstName: "Andi",
			LastName:  "Wijaya",
			Phone:     "+628133333333",
			Email:     "andi.wijaya@gmail.com",
		},
	}

	if err := db.Create(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}
