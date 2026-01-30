package data

import (
	"fmt"
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func ReservationSeeds(db *gorm.DB, customers []entity.Customer, tables []entity.Table) error {
	if len(customers) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	if len(tables) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	var count int64
	db.Model(&entity.Reservation{}).Count(&count)
	if count > 0 {
		return nil
	}
	
	reservations := []entity.Reservation{
		{
			CustomerID:      customers[0].ID,
			TableID:         tables[0].ID,
			PaxNumber:       4,
			ReservationDate: time.Now().AddDate(0, 0, 1), // besok
			ReservationTime: time.Date(0, 1, 1, 18, 30, 0, 0, time.Local),
			DepositFee:      50000,
			Status:          entity.ReservationStatusConfirmed,
			Notes:           "Reservasi ulang tahun",
		},
		{
			CustomerID:      customers[1].ID,
			TableID:         tables[1].ID,
			PaxNumber:       4,
			ReservationDate: time.Now().AddDate(0, 0, 1), // besok
			ReservationTime: time.Date(0, 1, 1, 18, 30, 0, 0, time.Local),
			DepositFee:      50000,
			Status:          entity.ReservationStatusConfirmed,
			Notes:           "Reservasi ulang tahun",
		},
	}

	return db.Create(&reservations).Error
}
