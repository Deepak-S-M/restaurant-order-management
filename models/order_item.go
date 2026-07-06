package models

type OrderItem struct {
	Model
	OrderID   string  `gorm:"type:text" json:"order_id"`
	ProductID string  `gorm:"type:text" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"type:numeric(10,2);not null" json:"unit_price"`
	Subtotal  float64 `gorm:"type:numeric(10,2);not null" json:"subtotal"`
}
