package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func InventoryLogSeeds(db *gorm.DB, products []entity.Product) error {
	if len(products) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	var count int64
	db.Model(&entity.Product{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	logs := []entity.InventoryLog{
		{
			ProductID:         products[0].ID,
			Type:              entity.InventoryLogTypeInitial,
			QuantityChange:    100,
			CurrentStockAfter: 100,
			Notes:             "Initial stock",
			CreatedBy:         1, // admin / system
		},
		{
			ProductID:         products[1].ID,
			Type:              entity.InventoryLogTypeInitial,
			QuantityChange:    80,
			CurrentStockAfter: 80,
			Notes:             "Initial stock",
			CreatedBy:         1,
		},
	}
	return db.Create(&logs).Error
}
