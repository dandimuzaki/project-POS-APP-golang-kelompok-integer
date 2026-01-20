package entity

import (
	"time"

	"gorm.io/gorm"
)

type TourSchedule struct {
	gorm.Model
	TourID    int
	StartDate time.Time
	EndDate time.Time
	PriceOverride float64
	Quota int
	BookedCount int
	Status string
}