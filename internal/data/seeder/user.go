package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"

	"gorm.io/gorm"
)

func UserSeeds(db *gorm.DB) ([]entity.User, error) {
	var count int64
	db.Model(&entity.User{}).Count(&count)

	var users []entity.User
	if count > 0 {
		db.Find(&users)
		return users, nil
	}

	users = []entity.User{
		{
			Email: "ihzhabaihaqqi05@gmail.com",
			PasswordHash: utils.HashPassword("password123"),
			Role: "superadmin",
		},
		{
			Email: "raflitbl1724@gmail.com",
			PasswordHash: utils.HashPassword("password123"),
			Role: "admin",
		},
		{
			Email: "dandimuzaki@gmail.com",
			PasswordHash: utils.HashPassword("password123"),
			Role: "staff",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}