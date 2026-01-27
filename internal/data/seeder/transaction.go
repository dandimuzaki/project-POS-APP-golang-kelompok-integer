package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

func TransactionSeeds() []entity.Transaction {
	return []entity.Transaction{
		{
			TransactionNumber: "TRX001",
			OrderID: 1,
			PaymentMethodID: 1,
			TransactionType: entity.TransactionTypePayment,
			Amount: 200000,
			Status: entity.TransactionStatusCompleted,
			CreatedBy: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			TransactionNumber: "TRX002",
			OrderID: 2,
			PaymentMethodID: 1,
			TransactionType: entity.TransactionTypePayment,
			Amount: 300000,
			Status: entity.TransactionStatusCompleted,
			CreatedBy: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}