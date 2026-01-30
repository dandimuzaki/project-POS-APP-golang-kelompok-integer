package data

import (
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
}