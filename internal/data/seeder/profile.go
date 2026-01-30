package data

import (
	"fmt"
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func StaffSeeds(db *gorm.DB, users []entity.User) error {
	if len(users) < 3 {
		return fmt.Errorf("not enough users to seed staff")
	}

	var count int64
	db.Model(&entity.Staff{}).Count(&count)
	if count > 0 {
		return nil
	}

	staffs := []entity.Staff{
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
	"project-POS-APP-golang-integer/internal/data/entity"
)

func ProfileSeeds() []entity.Profile{
	return []entity.Profile{
		{
			UserID: 3,
			FullName: "Rafli Nur Rahman",
			Phone: "087876544692",
			DateOfBirth: nil,
			Salary: 15000000,
			ProfileImageURL: "",
			Address: "Jl. Soekarno-Hatta, Tangerang, Banten",
			AdditionalDetails: "",
		},
		{
			UserID: 4,
			FullName: "Dandi Muhamad Zaki",
			Phone: "085117388153",
			DateOfBirth: nil,
			Salary: 15000000,
			ProfileImageURL: "",
			Address: "Jl. Gunung Batu, Bandung, Jawa Barat",
			AdditionalDetails: "",
		},
	}

	return db.Create(&staffs).Error
}
