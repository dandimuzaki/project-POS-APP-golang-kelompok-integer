package entity

import "gorm.io/gorm"

type Tour struct {
	gorm.Model
	LocationID  uint     `gorm:"not null" json:"location_id"`
	Name        string  `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string  `gorm:"type:varchar(200);uniqueIndex" json:"slug"`
	Description string  `gorm:"type:text" json:"description"`
	DurationDay int     `json:"duration_day"`
	BasePrice   float64 `gorm:"type:decimal(10,2);not null" json:"base_price"`
	Location Location `gorm:"foreignKey:LocationID" json:"location"`
	Category []*Category `gorm:"many2many:tour_categories" json:"categories"`
	Images []TourImage `gorm:"foreignKey:TourID" json:"images"`
	Itineraries []Itinerary `gorm:"foreignKey:TourID" json:"itineraries"`
	Schedules   []TourSchedule `gorm:"foreignKey:TourID" json:"schedules"`
	Reviews     []Review       `gorm:"foreignKey:TourID" json:"reviews"`
}

type TourImage struct {
	gorm.Model
	TourID uint `gorm:"not null" json:"tour_id"`
	ImageURL string `gorm:"type:text;not null" json:"image_url"`
	IsMain bool `gorm:"default:false" json:"is_main"`
}

type Itinerary struct {
	gorm.Model
	TourID uint `gorm:"not null" json:"tour_id"`
	DayNumber int `gorm:"not null" json:"day_number"`
	Title string `gorm:"type:varchar(150)" json:"title"`
	Description string `gorm:"type:text" json:"description"`
}