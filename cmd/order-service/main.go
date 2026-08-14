package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	"restaurant-order-management/middlewares"
	"restaurant-order-management/pb"
	"restaurant-order-management/seeder"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC in %s: %v\n%s\n", info.FullMethod, r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error: %v", r)
		}
	}()
	return handler(ctx, req)
}

func main() {
	fmt.Println("Order Service Starting...")
	config.ConnectDB(".env")
	// config.ConnectRedis()
	seeder.SeedData()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	api := r.Group("/api")
	// api.Use(middlewares.LoggerMiddleware())
	api.GET("/orders/:id", controllers.GetOrder)

	{
		// Protected routes
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Category
			protected.POST("/categories", middlewares.RoleMiddleware("admin"), controllers.CreateCategory)
			protected.GET("/categories/:id", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetCategory)
			protected.GET("/categories", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetCategories)
			protected.PUT("/categories/:id", middlewares.RoleMiddleware("admin"), controllers.UpdateCategory)
			protected.DELETE("/categories/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteCategory)

			// Products
			protected.POST("/products", middlewares.RoleMiddleware("admin"), controllers.CreateProduct)
			protected.GET("/products/:id", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetProduct)
			protected.GET("/products", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetProducts)
			protected.PUT("/products/:id", middlewares.RoleMiddleware("admin"), controllers.UpdateProduct)
			protected.DELETE("/products/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteProduct)

			// Orders
			protected.POST("/orders", middlewares.RoleMiddleware("admin", "waiter"), controllers.CreateOrder)
			// protected.GET("/orders/:id", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetOrder)
			protected.GET("/orders", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetOrders)
			protected.PUT("/orders/:id/status", middlewares.RoleMiddleware("admin", "waiter"), controllers.UpdateOrderStatus)
			protected.DELETE("/orders/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteOrder)
		}
	}

	addr := os.Getenv("ORDER_SERVICE_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	userServiceAddr := os.Getenv("USER_SERVICE_ADDRESS")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:8081"
	}

	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error connecting to User Service:", err)
		return
	}
	defer userConn.Close()

	userClient := pb.NewUserServiceClient(userConn)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting tcp")
		return
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(recoveryInterceptor),
	)

	pb.RegisterOrderServiceServer(grpcServer, &controllers.GrpcOrderServer{
		UserClient: userClient,
	})

	go func() {
		fmt.Println("Order Grpc server running on port: ", addr)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("Failed to start order grpc server...")
		}
	}()

	// srv := &http.Server{
	// 	Addr:    "localhost:8083",
	// 	Handler: r,
	// }

	// go func() {
	// 	fmt.Printf("Order Service running on %s\n", addr)
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatal("Failed to start server: ", err)
	// 	}
	// }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down Order Service...")
}
