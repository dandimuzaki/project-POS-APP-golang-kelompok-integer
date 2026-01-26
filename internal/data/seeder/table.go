package data

import "project-POS-APP-golang-integer/internal/data/entity"

func TableSeeds() []entity.Table {
	return []entity.Table{
		{
			TableNumber: "T01",
			Capacity:    2,
			Status:      entity.TableStatusAvailable,
		},
		{
			TableNumber: "T02",
			Capacity:    4,
			Status:      entity.TableStatusAvailable,
		},
		{
			TableNumber: "T03",
			Capacity:    4,
			Status:      entity.TableStatusAvailable,
		},
		{
			TableNumber: "T04",
			Capacity:    6,
			Status:      entity.TableStatusAvailable,
		},
		{
			TableNumber: "VIP-01",
			Capacity:    8,
			Status:      entity.TableStatusAvailable,
		},
	}
}
