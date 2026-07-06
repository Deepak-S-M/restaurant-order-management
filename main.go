package main

import (
	"fmt"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
)

func seedRoles() {

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
		}
		fmt.Println("Role seeded: ", role.Name)
	}

}

func main() {
	fmt.Println("Restuarant Order Management API Starting...")
	config.ConnectDB()
	seedRoles()
	fmt.Println("Server running on port 8080")
}
