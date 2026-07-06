package main

import (
	"fmt"
	"log"
	"restaurant-order-management/config"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Restuarant Order Management API Starting...")
	config.ConnectDB()
	SeedData()

	r := gin.Default()
	SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	} else {
		fmt.Println("Server running on port 8080")
	}
}
