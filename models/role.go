package models

type Role struct {
	Model
	Name string `gorm:"unique;not null"`
}
