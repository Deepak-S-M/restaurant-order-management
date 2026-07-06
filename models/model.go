package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	Id        string         `gorm:"primaryKey" json:"id" db:"id" uri:"id"`
	CreatedAt time.Time      `gorm:"index" json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `gorm:"index" json:"updated_at" db:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" db:"deleted_at"`
}

func (model *Model) BeforeCreate(tx *gorm.DB) (err error) {
	model.Id = uuid.New().String()
	return
}

func formatPreload(query *gorm.DB, preload ...string) {
	if len(preload) > 0 {
		for i := 0; i < len(preload); i++ {
			query.Preload(preload[i])
		}
	}
}
