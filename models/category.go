package models

type Category struct {
	Model
	Name        string    `gorm:"not null;type:text" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
