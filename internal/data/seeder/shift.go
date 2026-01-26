package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"time"
)

func ShiftSeeds() []entity.Shift{
	start := "09.00"
	end := "17.00"
	loc, _ := time.LoadLocation("Asia/Jakarta")

	wibStart, _ := time.ParseInLocation(
		"15.04",
		start,
		loc,
	)

	wibEnd, _ := time.ParseInLocation(
		"15.04",
		end,
		loc,
	)

	utcStart := wibStart.UTC()
	utcEnd := wibEnd.UTC()

	return []entity.Shift{
		{
			StaffID: 1,
			WeekNumber: 1,
			ShiftStart: utcStart,
			ShiftEnd: utcEnd,
			Year: 2026,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			StaffID: 2,
			WeekNumber: 2,
			ShiftStart: utcStart,
			ShiftEnd: utcEnd,
			Year: 2026,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}