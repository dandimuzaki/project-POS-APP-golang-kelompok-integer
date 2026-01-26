package data

import (
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func SeedSessions(db *gorm.DB) error {
	now := time.Now()
	exp := now.Add(24 * time.Hour)

	sessions := []entity.Session{
		{
			UserID:    1,
			Token:     "token_dummy_user_1",
			ExpiresAt: exp,
			CreatedAt: now,
		},
		{
			UserID:    2,
			Token:     "token_dummy_user_2",
			ExpiresAt: exp,
			CreatedAt: now,
		},
	}

	return db.Create(&sessions).Error
}
