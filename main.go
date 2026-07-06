package main

import (
	"fmt"
	"restaurant-order-management/config"
)

func main() {
	fmt.Println("Restuarant Order Management API Starting...")
	config.ConnectDB()
	fmt.Println("Server running on port 8080")
}
