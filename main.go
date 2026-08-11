package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"restaurant-order-management/handler"
	"restaurant-order-management/pb"
)

var (
	sseClients = make(map[string][]chan *pb.Order)
	sseMutex   sync.Mutex
)

func main() {

	orderConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()

	userConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer userConn.Close()

	orderHandler := handler.NewOrderServiceHandler(orderConn)
	userHandler := handler.NewUserServiceHandler(userConn)

	router := gin.Default()
	api := router.Group("/api")
	api.GET("/orders/:id", orderHandler.GetOrder)
	api.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)

	api.GET("/users/:id", userHandler.GetUser)

	log.Println("API Gateway running on :8080")

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
