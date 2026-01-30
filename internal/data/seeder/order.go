package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func OrderSeeds(db *gorm.DB, tables []entity.Table) ([]entity.Order, error) {
	var count int64
	db.Model(&entity.Order{}).Count(&count)

	var orders []entity.Order
	if count > 0 {
		db.Find(&orders)
		return orders, nil
	}

	orders = []entity.Order{
		{
			OrderNumber:   "ORD-001",
			TableID:       tables[0].ID,
			Status:        entity.OrderStatusCompleted,
			Subtotal:      50000,
			TaxPercentage: 10,
			TaxAmount:     5000,
			Total:         55000,
			CreatedBy:     1, // superadmin / admin
			Notes:         "Demo order",
		},
		{
			OrderNumber:   "ORD-002",
			TableID:       tables[1].ID,
			Status:        entity.OrderStatusCompleted,
			Subtotal:      50000,
			TaxPercentage: 10,
			TaxAmount:     5000,
			Total:         55000,
			CreatedBy:     1, // superadmin / admin
			Notes:         "Demo order",
		},
	}

	if err := db.Create(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
