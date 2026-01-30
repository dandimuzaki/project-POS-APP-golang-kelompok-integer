package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func OrderItemSeeds(db *gorm.DB, products []entity.Product, orders []entity.Order) error {
	if len(products) < 2 {
		return fmt.Errorf("not enough products to seed order_items")
	}

	var count int64
	db.Model(&entity.OrderItem{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	orderItems := []entity.OrderItem{
		{
			OrderID:    orders[0].ID,
			ProductID:  products[0].ID,
			Quantity:   2,
			TotalPrice: products[0].Price * 2,
		},
		{
			OrderID:    orders[0].ID,
			ProductID:  products[1].ID,
			Quantity:   3,
			TotalPrice: products[0].Price * 3,
		},
	}

	return db.Create(&orderItems).Error
}
