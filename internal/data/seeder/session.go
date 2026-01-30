package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func SessionSeeds(db *gorm.DB, users []entity.User) error {
	if len(users) < 3 {
		return fmt.Errorf("not enough users to seed sessions")
	}

	var count int64
	db.Model(&entity.Profile{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	token1, _ := utils.GenerateRandomToken(16)
	token2, _ := utils.GenerateRandomToken(16)
	token3, _ := utils.GenerateRandomToken(16)
	sessions := []entity.Session{
		{
			UserID: users[0].ID,
			Token: token1,
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
		{
			UserID: users[1].ID,
			Token: token2,
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
		{
			UserID: users[2].ID,
			Token: token3,
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
	}

	return db.Create(&sessions).Error
}