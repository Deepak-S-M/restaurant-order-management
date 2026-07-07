package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var roles []models.Role
	
	if err := config.DB.Find(&roles).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not fetch roles")
		return
	}

	utils.Success(c, http.StatusOK, "Roles fetched successfully", roles)
}

func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := config.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Role not found")
		return
	}

	utils.Success(c, http.StatusOK, "Role fetched successfully", role)
}
