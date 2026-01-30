package utils

import (
	"strconv"
	"strings"
	"time"
)

func BuildShiftTime(
	year int,
	week int,
	timeStr string,
	loc *time.Location,
) time.Time {
	parts := strings.Split(timeStr, ".")
	hour, _ := strconv.Atoi(parts[0])
	min, _ := strconv.Atoi(parts[1])

	// ISO week starts on Monday
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	for firstDay.Weekday() != time.Monday {
		firstDay = firstDay.AddDate(0, 0, 1)
	}

	shiftDate := firstDay.AddDate(0, 0, (week-1)*7)
	return time.Date(
		shiftDate.Year(),
		shiftDate.Month(),
		shiftDate.Day(),
		hour,
		min,
		0,
		0,
		loc,
	).UTC()
}
