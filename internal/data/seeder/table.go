package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func TableSeeds (db *gorm.DB) ([]entity.Table, error) {
	var count int64
	db.Model(&entity.Table{}).Count(&count)

	var tables []entity.Table
	if count > 0 {
		db.Find(&tables)
		return tables, nil
	}

	tables = []entity.Table{
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

	if err := db.Create(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}
