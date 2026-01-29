package entity

import (
	"time"
)

type Shift struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StaffID    uint      `gorm:"not null" json:"staff_id"`
	WeekNumber int       `gorm:"not null" json:"week_number"`
	ShiftStart time.Time `json:"shift_start"`
	ShiftEnd   time.Time `json:"shift_end"`
	Year       int       `gorm:"not null" json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Profile Profile `gorm:"foreignKey:StaffID" json:"profile,omitempty"`
}
