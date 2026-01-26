package data

import (
	"fmt"
	"time"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/gorm"
)

func SeedShifts(db *gorm.DB, staffs []entity.Staff) error {
	if len(staffs) == 0 {
		return fmt.Errorf("no staff found for shift seeding")
	}

	var count int64
	db.Model(&entity.Shift{}).Count(&count)
	if count > 0 {
		// shift sudah ada → skip
		return nil
	}

	currentYear := time.Now().Year()

	shifts := []entity.Shift{}

	for _, staff := range staffs {
		// contoh: 2 shift per staff
		shifts = append(shifts,
			entity.Shift{
				StaffID:    staff.ID,
				WeekNumber: 1,
				ShiftStart: time.Date(currentYear, 1, 1, 8, 0, 0, 0, time.UTC),
				ShiftEnd:   time.Date(currentYear, 1, 1, 16, 0, 0, 0, time.UTC),
				Year:       currentYear,
			},
			entity.Shift{
				StaffID:    staff.ID,
				WeekNumber: 1,
				ShiftStart: time.Date(currentYear, 1, 2, 16, 0, 0, 0, time.UTC),
				ShiftEnd:   time.Date(currentYear, 1, 2, 23, 0, 0, 0, time.UTC),
				Year:       currentYear,
			},
		)
	}

	return db.Create(&shifts).Error
}
