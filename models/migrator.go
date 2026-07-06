package models

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&Role{},
		&User{},
		&Category{},
		&Product{},
		&Order{},
		&OrderItem{},
	); err != nil {
		log.Println("Error migrating tables: ", err)
	}
}
