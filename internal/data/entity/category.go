package entity

type Category struct {
	Model
	Name  string  `json:"name"`
	Tours []*Tour `gorm:"many2many:tour_categories;" json:"tours,omitempty"`
}