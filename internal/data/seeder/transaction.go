package data

import "project-POS-APP-golang-integer/internal/data/entity"

func TransactionSeeds() []entity.Transaction {
	return []entity.Transaction{
		{
			TransactionNumber: "TRX-0001",
			OrderID:           1,
			TransactionType:   entity.TransactionTypePayment,
			PaymentMethodID:   1,
			Amount:            150000,
			Status:            entity.TransactionStatusCompleted,
			CreatedBy:         1,
		},
	}
}
