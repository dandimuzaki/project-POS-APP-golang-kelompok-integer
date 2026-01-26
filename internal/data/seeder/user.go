package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
)

func UserSeeds() []entity.User {
	return []entity.User{
		{
			Email: "integer.lumoshive@gmail.com",
			PasswordHash: utils.HashPassword("password123"),
			Role: "superadmin",
		},
		{
			Email: "dandimuzaki@gmail.com",
			PasswordHash: utils.HashPassword("password123"),
			Role: "admin",
		},
	}
}