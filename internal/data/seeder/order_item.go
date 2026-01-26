package data

import "project-POS-APP-golang-integer/internal/data/entity"

func OrderItemSeeds() []entity.OrderItem {
	return []entity.OrderItem{
		{
			OrderID:    1,
			ProductID:  1,
			Quantity:   2,
			TotalPrice: 50000,
		},
	}
}
