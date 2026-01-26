package data

import "project-POS-APP-golang-integer/internal/data/entity"

func PaymentMethodSeeds() []entity.PaymentMethod {
	return []entity.PaymentMethod{
		{
			Name:     "Cash",
			IsActive: true,
		},
		{
			Name:     "Debit Card",
			IsActive: true,
		},
		{
			Name:     "Credit Card",
			IsActive: true,
		},
		{
			Name:     "QRIS",
			IsActive: true,
		},
		{
			Name:     "E-Wallet",
			IsActive: true,
		},
	}
}
