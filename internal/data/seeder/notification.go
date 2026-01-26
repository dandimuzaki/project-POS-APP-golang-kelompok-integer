package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

func NotificationSeeds() []entity.Notification {
	return []entity.Notification{
		{
			UserID: 1,
			Title: "Welcome",
			Message: "Welcome to POS System created by Integer",
			Type: entity.NotificationTypeSystem,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID: 2,
			Title: "Welcome",
			Message: "Welcome to POS System created by Integer",
			Type: entity.NotificationTypeSystem,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID: 3,
			Title: "Welcome",
			Message: "Welcome to POS System created by Integer",
			Type: entity.NotificationTypeSystem,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID: 4,
			Title: "Welcome",
			Message: "Welcome to POS System created by Integer",
			Type: entity.NotificationTypeSystem,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}