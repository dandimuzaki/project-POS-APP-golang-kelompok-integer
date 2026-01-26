package data

import "project-POS-APP-golang-integer/internal/data/entity"

func NotificationSeeds() []entity.Notification {
	return []entity.Notification{
		{
			UserID:  1,
			Title:   "Welcome",
			Message: "Selamat datang di POS System. Akun Anda telah aktif.",
			Type:    entity.NotificationTypeSystem,
			Status:  entity.NotificationStatusNew,
		},
		{
			UserID:  2,
			Title:   "System Info",
			Message: "Silakan lengkapi data profil Anda.",
			Type:    entity.NotificationTypeSystem,
			Status:  entity.NotificationStatusNew,
		},
	}
}
