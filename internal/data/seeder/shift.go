package data

import (
	"fmt"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ShiftSeeds(db *gorm.DB, profiles []entity.Profile) error {
	if len(profiles) < 2 {
		return fmt.Errorf("not enough users to seed profile")
	}

	var count int64
	db.Model(&entity.Shift{}).Count(&count)
	if count > 0 {
		return nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	startTime := "09.00"
	endTime := "17.00"

	var shifts []entity.Shift

	for i, profile := range profiles {
		week := i + 1

		shiftStart := utils.BuildShiftTime(2026, week, startTime, loc)
		shiftEnd := utils.BuildShiftTime(2026, week, endTime, loc)

		shifts = append(shifts, entity.Shift{
			ProfileID:  profile.ID,
			WeekNumber: week,
			Year:       2026,
			ShiftStart: shiftStart,
			ShiftEnd:   shiftEnd,
		})
	}

	return db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "profile_id"},
				{Name: "week_number"},
				{Name: "year"},
			},
			DoNothing: true,
		}).
		Create(&shifts).Error
}
