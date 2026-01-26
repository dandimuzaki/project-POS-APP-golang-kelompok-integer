package data

import "project-POS-APP-golang-integer/internal/data/entity"

func InventoryLogSeeds() []entity.InventoryLog {
	return []entity.InventoryLog{
		{
			ProductID:         1,
			Type:              entity.InventoryLogTypeInitial,
			QuantityChange:    100,
			CurrentStockAfter: 100,
			Notes:             "Initial stock",
			CreatedBy:         1, // admin / system
		},
		{
			ProductID:         2,
			Type:              entity.InventoryLogTypeInitial,
			QuantityChange:    80,
			CurrentStockAfter: 80,
			Notes:             "Initial stock",
			CreatedBy:         1,
		},
	}
}
