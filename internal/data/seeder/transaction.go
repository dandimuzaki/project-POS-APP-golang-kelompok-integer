package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"

	"gorm.io/gorm"
)

func TransactionSeeds(db *gorm.DB, orders []entity.Order, pMethods []entity.PaymentMethod) error {
	if len(orders) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	if len(pMethods) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	var count int64
	db.Model(&entity.Transaction{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	transactions := []entity.Transaction{
		{
			TransactionNumber: "TRX001",
			OrderID: orders[0].ID,
			PaymentMethodID: pMethods[0].ID,
			TransactionType: entity.TransactionTypePayment,
			Amount: 200000,
			Status: entity.TransactionStatusCompleted,
			CreatedBy: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			TransactionNumber: "TRX002",
			OrderID: orders[1].ID,
			PaymentMethodID: pMethods[0].ID,
			TransactionType: entity.TransactionTypePayment,
			Amount: 300000,
			Status: entity.TransactionStatusCompleted,
			CreatedBy: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return db.Create(&transactions).Error
}