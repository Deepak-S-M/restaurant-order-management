package main

import (
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	"restaurant-order-management/handler"
	"restaurant-order-management/middlewares"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	config.ConnectDB()

	orderConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()

	userConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()

	orderHandler := handler.NewOrderServiceHandler(orderConn)
	userHandler := handler.NewUserServiceHandler(userConn)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(middlewares.LoggerMiddleware())
	api.POST("/login", controllers.Login)
	{
		protectedApi := api.Group("/")
		protectedApi.Use(middlewares.AuthMiddleware())
		{
			protectedApi.GET("/users/:id", middlewares.RoleMiddleware("admin"), userHandler.GetUser)
			protectedApi.GET("/order/:id", middlewares.RoleMiddleware("admin", "waiter"), orderHandler.GetOrder)
			protectedApi.PUT("/orders/:id/status", middlewares.RoleMiddleware("admin", "waiter"), orderHandler.UpdateOrderStatus)
		}
	}

	log.Println("API Gateway running on :8000")

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
