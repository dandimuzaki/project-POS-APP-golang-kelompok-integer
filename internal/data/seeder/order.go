package data

import "project-POS-APP-golang-integer/internal/data/entity"

func OrderSeeds() []entity.Order {
	return []entity.Order{
		{
			OrderNumber:   "ORD-001",
			TableID:       1,
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
			TableID:       2,
			Status:        entity.OrderStatusCompleted,
			Subtotal:      50000,
			TaxPercentage: 10,
			TaxAmount:     5000,
			Total:         55000,
			CreatedBy:     1, // superadmin / admin
			Notes:         "Demo order",
		},
	}
}
