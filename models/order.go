package models

type Order struct {
	Model
	UserID     string      `gorm:"type:text" json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status     string      `gorm:"type:text;default:'pending'" json:"status"`
	Subtotal   float64     `gorm:"type:numeric(10,2);not null" json:"subtotal"`
	Tax        float64     `gorm:"type:numeric(10,2);not null" json:"tax"`
	GrandTotal float64     `gorm:"type:numeric(10,2);not null" json:"grand_total"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}
