package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Auth
		&entity.User{}, 
		&entity.Profile{}, 
		&entity.Shift{},
		&entity.OTP{},
		&entity.Session{},

		// Menu
		&entity.Category{},
		&entity.Product{},
		&entity.InventoryLog{},
		
		// Order
		&entity.Table{},
		&entity.Customer{},
		&entity.Order{},
		&entity.Reservation{},
		&entity.OrderItem{},
		
		// Payment
		&entity.PaymentMethod{},
		&entity.Transaction{},

		&entity.Notification{},
	)
}
