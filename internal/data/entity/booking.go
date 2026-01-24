package entity

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserID         uint         `gorm:"not null" json:"user_id"`
	TourScheduleID uint         `gorm:"not null" json:"tour_schedule_id"`
	BookingCode    string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"booking_code"`
	NumPeople      int          `gorm:"not null" json:"num_people"`
	TotalAmount    float64      `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status         string       `gorm:"type:varchar(20);default:'pending'" json:"status"`
	User           User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TourSchedule   TourSchedule `gorm:"foreignKey:TourScheduleID" json:"tour_schedule,omitempty"`
	Payment        *Payment     `gorm:"foreignKey:BookingID" json:"payment,omitempty"`
}
