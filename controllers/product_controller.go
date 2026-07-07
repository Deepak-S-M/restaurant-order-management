package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductInput struct {
	CategoryID  *string  `json:"category_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	ImageURL    *string  `json:"image_url"`
}

func CreateProduct(c *gin.Context) {
	var input ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product := models.Product{
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		ImageURL:    input.ImageURL,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not create product")
		return
	}

	utils.Success(c, http.StatusCreated, "Product created successfully", product)
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	db := config.DB.Model(&models.Product{})

	// Filtering by category
	categoryID := c.Query("category_id")
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// Searching by name
	search := c.Query("search")
	if search != "" {
		db = db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	db.Count(&total)

	if err := db.Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not fetch products")
		return
	}

	utils.SuccessWithPagination(c, http.StatusOK, "Products fetched successfully", products, total, int64(page), int64(limit))
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.Preload("Category").Where("id = ?", id).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	utils.Success(c, http.StatusOK, "Product fetched successfully", product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Stock != nil {
		product.Stock = *input.Stock
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}

	if err := config.DB.Save(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not update product")
		return
	}

	utils.Success(c, http.StatusOK, "Product updated successfully", product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not delete product")
		return
	}

	utils.Success(c, http.StatusOK, "Product deleted successfully", nil)
}
