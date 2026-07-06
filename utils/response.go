package utils

import (
	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Errors     any         `json:"errors,omitempty"`
}

func Success(ctx *gin.Context, statusCode int, message string, data any) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(ctx *gin.Context, statusCode int, message string, data any, total int64, page int64, pageSize int64) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: &Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}

func ValidationError(ctx *gin.Context, statusCode int, message string, errors any) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
