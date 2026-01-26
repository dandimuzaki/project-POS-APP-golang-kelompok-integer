package data

import (
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"
)

func ReservationSeeds() []entity.Reservation {
	return []entity.Reservation{
		{
			CustomerID:      1,
			TableID:         2,
			PaxNumber:       4,
			ReservationDate: time.Now().AddDate(0, 0, 1), // besok
			ReservationTime: time.Date(0, 1, 1, 18, 30, 0, 0, time.Local),
			DepositFee:      50000,
			Status:          entity.ReservationStatusConfirmed,
			Notes:           "Reservasi ulang tahun",
		},
	}
}
