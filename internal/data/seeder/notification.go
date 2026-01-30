package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func NotificationSeeds(db *gorm.DB, users []entity.User) error {
	if len(users) < 3 {
		return fmt.Errorf("not enough users to seed profile")
	}

	var count int64
	db.Model(&entity.Profile{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	notifications := []entity.Notification{
		{
			UserID:  users[0].ID,
			Title:   "Welcome",
			Message: "Selamat datang di POS System. Akun Anda telah aktif.",
			Type:    entity.NotificationTypeSystem,
			Status:  entity.NotificationStatusNew,
		},
		{
			UserID:  users[1].ID,
			Title:   "System Info",
			Message: "Silakan lengkapi data profil Anda.",
			Type:    entity.NotificationTypeSystem,
			Status:  entity.NotificationStatusNew,
		},
	}
	return db.Create(&notifications).Error
}
