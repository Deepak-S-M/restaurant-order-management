package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func CreateCategory(c *gin.Context) {
	var input CategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not create category")
		return
	}

	utils.Success(c, http.StatusCreated, "Category created successfully", category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := config.DB.Find(&categories).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not fetch categories")
		return
	}

	utils.Success(c, http.StatusOK, "Categories fetched successfully", categories)
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.DB.Where("id = ?", id).First(&category).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.Success(c, http.StatusOK, "Category fetched successfully", category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.DB.Where("id = ?", id).First(&category).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}

	if err := config.DB.Save(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not update category")
		return
	}

	utils.Success(c, http.StatusOK, "Category updated successfully", category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := config.DB.Where("id = ?", id).First(&category).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not delete category")
		return
	}

	utils.Success(c, http.StatusOK, "Category deleted successfully", nil)
}
