package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	// Fetch user with Role
	if err := config.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Verify Password
	if err := utils.CheckPassword(input.Password, user.Password); err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate Token
	token, err := utils.GenerateToken(user.Id, user.Email, user.Role.Name)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not generate token")
		return
	}

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role.Name,
		},
	})
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, "Email already in use")
		return
	}

	roleName := "waiter"
	if input.Role != "" {
		roleName = input.Role
	}

	var role models.Role
	if err := config.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid role")
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

	utils.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": gin.H{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"role":  role.Name,
		},
	})
}
