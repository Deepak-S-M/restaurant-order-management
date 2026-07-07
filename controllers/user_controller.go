package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"

	"github.com/gin-gonic/gin"
)

type UpdateUserInput struct {
	Name   *string `json:"name"`
	Email  *string `json:"email" binding:"omitempty,email"`
	RoleID *string `json:"role_id"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, "Email already in use")
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ?", input.RoleID).First(&role).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid Role ID")
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not hash password")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		RoleId:   role.Id,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not create user")
		return
	}

	user.Password = ""
	utils.Success(c, http.StatusCreated, "User created successfully", user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not fetch users")
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	utils.Success(c, http.StatusOK, "Users fetched successfully", users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	user.Password = ""
	utils.Success(c, http.StatusOK, "User fetched successfully", user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Email != nil {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", *input.Email, id).First(&existingUser).Error; err == nil {
			utils.Error(c, http.StatusConflict, "Email already in use")
			return
		}
		user.Email = *input.Email
	}

	if input.RoleID != nil {
		var role models.Role
		if err := config.DB.Where("id = ?", *input.RoleID).First(&role).Error; err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid Role ID")
			return
		}
		user.RoleId = *input.RoleID
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not update user")
		return
	}

	user.Password = ""
	utils.Success(c, http.StatusOK, "User updated successfully", user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not delete user")
		return
	}

	utils.Success(c, http.StatusOK, "User deleted successfully", nil)
}
