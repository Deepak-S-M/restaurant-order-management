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

func SeedCategories() {
	var count int64
	config.DB.Model(&models.Category{}).Count(&count)

	if count > 0 {
		fmt.Println("Categories already seeded, skipping...")
		return
	}

	categories := []models.Category{
		{Name: "Appetizers", Description: "Starters and light snacks"},
		{Name: "Main Course", Description: "Primary dishes"},
		{Name: "Desserts", Description: "Sweet treats"},
		{Name: "Beverages", Description: "Drinks and refreshments"},
	}

	for _, category := range categories {
		result := config.DB.Create(&category)
		if result.Error != nil {
			log.Println("Error in seeding category: ", result.Error)
		} else {
			fmt.Println("Category seeded: ", category.Name)
		}
	}
}

func SeedProducts() {
	var count int64
	config.DB.Model(&models.Product{}).Count(&count)

	if count > 0 {
		fmt.Println("Products already seeded, skipping...")
		return
	}

	var mainCourse models.Category
	if err := config.DB.Where("name = ?", "Main Course").First(&mainCourse).Error; err != nil {
		log.Println("Main Course category not found, cannot seed products")
		return
	}

	products := []models.Product{
		{
			CategoryID:  mainCourse.Id,
			Name:        "Classic Cheeseburger",
			Description: "Beef patty with cheddar cheese, lettuce, and tomato",
			Price:       12.99,
			Stock:       50,
		},
		{
			CategoryID:  mainCourse.Id,
			Name:        "Margherita Pizza",
			Description: "Classic pizza with tomato sauce, mozzarella, and basil",
			Price:       14.50,
			Stock:       30,
		},
	}

	for _, product := range products {
		result := config.DB.Create(&product)
		if result.Error != nil {
			log.Println("Error in seeding product: ", result.Error)
		} else {
			fmt.Println("Product seeded: ", product.Name)
		}
	}
}

func SeedData() {
	SeedRoles()
	SeedUsers()
	SeedCategories()
	SeedProducts()
}
