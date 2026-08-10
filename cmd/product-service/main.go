package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	"restaurant-order-management/middlewares"
	"restaurant-order-management/seeder"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Product Service Starting...")
	config.ConnectDB()
	config.ConnectRedis()
	seeder.SeedData()

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	api.Use(middlewares.LoggerMiddleware())
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
		}
	}

	addr := os.Getenv("PRODUCT_SERVICE_ADDRESS")
	if addr == "" {
		addr = ":8082"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		fmt.Printf("Product Service running on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down Product Service...")
}
