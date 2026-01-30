package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"

	"gorm.io/gorm"
)

func PaymentMethodSeeds(db *gorm.DB) ([]entity.PaymentMethod, error) {
	var count int64
	db.Model(&entity.PaymentMethod{}).Count(&count)

	var pMethods []entity.PaymentMethod
	if count > 0 {
		db.Find(&pMethods)
		return pMethods, nil
	}

	pMethods = []entity.PaymentMethod{
		{
			Name: "Cash",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name: "Debit Card",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name: "E-Wallet",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&pMethods).Error; err != nil {
		return nil, err
	}

	return pMethods, nil
}