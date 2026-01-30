package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func OTPSeeds(db *gorm.DB, users []entity.User) error {
	if len(users) < 3 {
		return fmt.Errorf("not enough users to seed otp")
	}

	var count int64
	db.Model(&entity.OTP{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	otp1, _ := utils.GenerateOTP(6)
	otp2, _ := utils.GenerateOTP(6)
	otp3, _ := utils.GenerateOTP(6)
	otps := []entity.OTP{
		{
			UserID: users[0].ID,
			OTPCode: otp1,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
		{
			UserID: users[1].ID,
			OTPCode: otp2,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
		{
			UserID: users[2].ID,
			OTPCode: otp3,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			IsUsed: false,
		},
	}
	return db.Create(&otps).Error
}