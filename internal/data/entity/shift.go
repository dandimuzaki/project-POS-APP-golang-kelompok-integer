package entity

import (
	"time"
)

type Shift struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProfileID    uint      `gorm:"not null" json:"profile_id"`
	WeekNumber int       `gorm:"not null" json:"week_number"`
	ShiftStart time.Time `json:"shift_start"`
	ShiftEnd   time.Time `json:"shift_end"`
	Year       int       `gorm:"not null" json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Profile Profile `gorm:"foreignKey:ProfileID" json:"profile,omitempty"`
}
