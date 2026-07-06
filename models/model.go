package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

type Pagination struct {
	Page     int64 `form:"page" binding:"gte=1"`
	PageSize int64 `form:"pageSize" binding:"gte=1,lte=100"`
	Offset   int
	Limit    int
	SortBy   string `form:"sortBy"`
	Total    int64
}

func BindQueryToPagination(ctx *gin.Context) (pagination *Pagination) {
	page := ctx.Query("page")
	pageSize := ctx.Query("pageSize")
	sortBy := ctx.Query("sortBy")

	if page != "" && pageSize != "" {
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			fmt.Println("Failed to convert")
		}
		pageSizeNum, err := strconv.Atoi(pageSize)
		if err != nil {
			fmt.Println("Failed to convert")
		}
		pagination := Pagination{
			Page:     int64(pageNum),
			PageSize: int64(pageSizeNum),
			SortBy:   sortBy,
			Offset:   (pageNum - 1) * pageSizeNum,
			Limit:    pageSizeNum,
		}
		return &pagination
	}
	return &Pagination{
		Page:     0,
		PageSize: 0,
		SortBy:   "",
		Offset:   0,
		Limit:    -1,
	}
}

func formatPagination(query *gorm.DB, pagination *Pagination, defaultSort string) {
	if pagination != nil {
		if pagination.SortBy == "" {
			pagination.SortBy = defaultSort
		}
		query.Offset(pagination.Offset).Order(pagination.SortBy).Limit(pagination.Limit)
	} else {
		query.Order(defaultSort)
	}
}
