package main

import (
	"restaurant-order-management/controllers"
	"restaurant-order-management/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middlewares.LoggerMiddleware())
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

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
			protected.GET("/orders/:id", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetOrder)
			protected.GET("/orders", middlewares.RoleMiddleware("admin", "waiter"), controllers.GetOrders)
			protected.PUT("/orders/:id/status", middlewares.RoleMiddleware("admin", "waiter"), controllers.UpdateOrderStatus)
			protected.DELETE("/orders/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteOrder)

			// Users (Admin only)
			protected.POST("/users", middlewares.RoleMiddleware("admin"), controllers.CreateUser)
			protected.GET("/users", middlewares.RoleMiddleware("admin"), controllers.GetUsers)
			protected.GET("/users/:id", middlewares.RoleMiddleware("admin"), controllers.GetUser)
			protected.PUT("/users/:id", middlewares.RoleMiddleware("admin"), controllers.UpdateUser)
			protected.DELETE("/users/:id", middlewares.RoleMiddleware("admin"), controllers.DeleteUser)
		}
	}
}
