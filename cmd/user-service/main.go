package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	"restaurant-order-management/middlewares"
	"restaurant-order-management/pb"
	"restaurant-order-management/seeder"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("User Service Starting...")
	config.ConnectDB(".env")
	config.ConnectRedis()
	seeder.SeedData()

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	api.Use(middlewares.LoggerMiddleware())
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Users (Admin only)
			protected.POST("/users", middlewares.RoleMiddleware("admin"), controllers.CreateUser)
			protected.GET("/users", middlewares.RoleMiddleware("admin"), controllers.GetUsers)
			protected.GET("/users/:id", middlewares.RoleMiddleware("admin"), controllers.GetUser)
			protected.PUT("/users/:id", middlewares.RoleMiddleware("admin"), controllers.UpdateUser)
			protected.DELETE("/users/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteUser)

			// Roles (Admin only)
			protected.GET("/roles", middlewares.RoleMiddleware("admin"), controllers.GetRoles)
			protected.GET("/roles/:id", middlewares.RoleMiddleware("admin"), controllers.GetRole)
		}
	}

	// addr := os.Getenv("USER_SERVICE_ADDRESS")
	// if addr == "" {
	// 	addr = ":8081"
	// }

	// srv := &http.Server{
	// 	Addr:    addr,
	// 	Handler: r,
	// }

	// go func() {
	// 	fmt.Printf("User Service running on %s\n", addr)
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatal("Failed to start server: ", err)
	// 	}
	// }()

	grpcAddr := os.Getenv("USER_SERVICE_GRPC_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = ":8081"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		fmt.Printf("Failed to listen on gRPC port: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &controllers.GrpcUserServer{})

	go func() {
		fmt.Printf("User Service gRPC running on %s\n", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("Failed to start gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down User Service...")
}
