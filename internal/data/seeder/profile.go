package data

import (
	"fmt"
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func ProfileSeeds(db *gorm.DB, users []entity.User) ([]entity.Profile, error) {
	if len(users) < 3 {
		return nil, fmt.Errorf("not enough users to seed profile")
	}

	var profiles []entity.Profile
	var count int64
	if count > 0 {
		db.Find(&profiles)
		return profiles, nil
	}

	profiles = []entity.Profile{
		{
			UserID:            users[0].ID,
			FullName:          "Super Admin",
			Phone:             "+628111111111",
			DateOfBirth:       time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC),
			Salary:            15_000_000,
			Address:           "Jakarta",
			AdditionalDetails: "System Administrator",
		},
		{
			UserID:            users[1].ID,
			FullName:          "Restaurant Manager",
			Phone:             "+628122222222",
			DateOfBirth:       time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC),
			Salary:            12_000_000,
			Address:           "Jakarta",
			AdditionalDetails: "Restaurant Manager",
		},
		{
			UserID:            users[2].ID,
			FullName:          "Waitress",
			Phone:             "+628133333333",
			DateOfBirth:       time.Date(1995, 11, 5, 0, 0, 0, 0, time.UTC),
			Salary:            7_500_000,
			Address:           "Jakarta",
			AdditionalDetails: "Senior Waitress",
		},
	}

	if err := db.Create(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}
