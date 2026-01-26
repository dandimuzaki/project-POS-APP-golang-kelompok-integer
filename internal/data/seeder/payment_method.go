package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

func PaymentMethodSeeds() []entity.PaymentMethod{
	return []entity.PaymentMethod{
		{
			Name: "Gopay",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name: "OVO",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name: "BCA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}