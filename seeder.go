package main

import (
	"fmt"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"
)

func SeedRoles() {
	var count int64
	result := config.DB.Model(&models.Role{}).Count(&count)
	if result.Error != nil {
		log.Println("Unable to count role: ", result.Error)
		return
	}

	if count > 0 {
		fmt.Println("Roles already seeded, skipping...")
		return
	}

	roles := []models.Role{
		{Name: "admin"},
		{Name: "waiter"},
	}

	for _, role := range roles {
		result := config.DB.Model(&models.Role{}).Create(&role)
		if result.Error != nil {
			log.Println("Error in seeding roles: ", result.Error)
		} else {
			fmt.Println("Role seeded: ", role.Name)
		}
	}
}

func SeedUsers() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		fmt.Println("Users already seeded, skipping...")
		return
	}

	var adminRole models.Role
	err := config.DB.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		log.Println("Admin role not found, make sure roles are seeded first")
		return
	}

	var waiterRole models.Role
	err = config.DB.Where("name = ?", "waiter").First(&waiterRole).Error
	if err != nil {
		log.Println("Waiter role not found, make sure roles are seeded first")
		return
	}

	adminPassword, _ := utils.HashPassword("admin123")
	waiterPassword, _ := utils.HashPassword("waiter123")

	users := []models.User{
		{
			Name:     "Admin User",
			Email:    "admin@restaurant.com",
			Password: adminPassword,
			RoleId:   adminRole.Id,
		},
		{
			Name:     "Waiter User",
			Email:    "waiter@restaurant.com",
			Password: waiterPassword,
			RoleId:   waiterRole.Id,
		},
	}

	for _, user := range users {
		result := config.DB.Create(&user)
		if result.Error != nil {
			log.Println("Error in seeding user: ", result.Error)
		} else {
			fmt.Println("User seeded: ", user.Name, "(Email:", user.Email, ")")
		}
	}
}

func SeedData() {
	SeedRoles()
	SeedUsers()
}
