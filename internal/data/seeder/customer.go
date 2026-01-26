package data

import "project-POS-APP-golang-integer/internal/data/entity"

func CustomerSeeds() []entity.Customer {
	return []entity.Customer{
		{
			Title: entity.CustomerTitleMr,
			FirstName: "Andi",
			LastName: "Sopandi",
			Phone: "082134657809",
			Email: "andi@gmail.com",
		},
		{
			Title: entity.CustomerTitleMr,
			FirstName: "Siti",
			LastName: "Rahmawati",
			Phone: "082134657899",
			Email: "siti@gmail.com",
		},
	}
}