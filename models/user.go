package models

type User struct {
	Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	RoleId   string `json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleId" json:"role"`
}
