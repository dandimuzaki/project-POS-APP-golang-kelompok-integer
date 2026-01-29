package entity

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID            uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	FullName          string    `gorm:"not null" json:"full_name"`
	Phone             string    `json:"phone"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	Salary            float64   `json:"salary"`
	ProfileImageURL   string    `json:"profile_image_url"`
	Address           string    `json:"address"`
	AdditionalDetails string    `gorm:"type:text" json:"additional_details,omitempty"`

	// Relations
	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Shifts []Shift `gorm:"foreignKey:StaffID" json:"shifts,omitempty"`
}
