package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"
)

func OTPSeeds() []entity.OTP {
	otp1, _ := utils.GenerateOTP(6)
	otp2, _ := utils.GenerateOTP(6)
	otp3, _ := utils.GenerateOTP(6)
	return []entity.OTP{
		{
			UserID: 1,
			OTPCode: otp1,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
		{
			UserID: 2,
			OTPCode: otp2,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
		{
			UserID: 3,
			OTPCode: otp3,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
	}
}