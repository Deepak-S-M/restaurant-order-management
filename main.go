package main

import (
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/handler"
	"restaurant-order-management/middlewares"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "restaurant-order-management/docs"
)

// @title Restaurant Order Management API
// @version 1.0
// @description Restaurant Order Management Microservices API
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", userHandler.Login)
	api := router.Group("/api")
	api.Use(middlewares.LoggerMiddleware())
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
