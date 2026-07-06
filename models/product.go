package models

type Product struct {
	Model
	CategoryID  string   `gorm:"type:text" json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string   `gorm:"not null;type:text" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	Price       float64  `gorm:"type:numeric(10,2);not null" json:"price"`
	Stock       int      `gorm:"default:0" json:"stock"`
	ImageURL    string   `gorm:"type:text" json:"image_url"`
}
