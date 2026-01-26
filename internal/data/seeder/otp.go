package data

import (
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"
)

func OTPSeeds() []entity.OTP {
	now := time.Now()

	return []entity.OTP{
		{
			UserID:    1,
			OTPCode:   "123456",
			ExpiresAt: now.Add(5 * time.Minute),
			IsUsed:    false,
		},
		{
			UserID:    2,
			OTPCode:   "654321",
			ExpiresAt: now.Add(5 * time.Minute),
			IsUsed:    false,
		},
	}
}
