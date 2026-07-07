package controllers

import (
	"net/http"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderItemInput struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderInput struct {
	Items []OrderItemInput `json:"items" binding:"required,min=1"`
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var subtotal float64
	var orderItems []models.OrderItem

	// Calculate totals and prepare items
	for _, itemInput := range input.Items {
		var product models.Product
		if err := config.DB.Where("id = ?", itemInput.ProductID).First(&product).Error; err != nil {
			utils.Error(c, http.StatusBadRequest, "Invalid product ID: "+itemInput.ProductID)
			return
		}

		if product.Stock < itemInput.Quantity {
			utils.Error(c, http.StatusBadRequest, "Insufficient stock for product: "+product.Name)
			return
		}

		itemSubtotal := product.Price * float64(itemInput.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.Id,
			Quantity:  itemInput.Quantity,
			UnitPrice: product.Price,
			Subtotal:  itemSubtotal,
		})
	}

	tax := subtotal * 0.10 // 10% tax rate
	grandTotal := subtotal + tax

	order := models.Order{
		UserID:     userID.(string),
		Status:     "pending",
		Subtotal:   subtotal,
		Tax:        tax,
		GrandTotal: grandTotal,
	}

	// Transaction to create order and items
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.Id
			if err := tx.Create(&orderItems[i]).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Product{}).Where("id = ?", orderItems[i].ProductID).UpdateColumn("stock", gorm.Expr("stock - ?", orderItems[i].Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not create order")
		return
	}

	// Fetch full order to return
	config.DB.Preload("Items.Product").Where("id = ?", order.Id).First(&order)

	utils.Success(c, http.StatusCreated, "Order created successfully", order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	db := config.DB.Model(&models.Order{})

	// Filter by user if not admin
	role, _ := c.Get("role")
	if role == "waiter" {
		userID, _ := c.Get("user_id")
		db = db.Where("user_id = ?", userID)
	}

	// Filter by status
	status := c.Query("status")
	if status != "" {
		db = db.Where("LOWER(status) = ?", strings.ToLower(status))
	}

	// Search by ID
	search := c.Query("search")
	if search != "" {
		db = db.Where("id = ?", search)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	db.Count(&total)

	if err := db.Preload("User").Preload("Items.Product").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not fetch orders")
		return
	}

	utils.SuccessWithPagination(c, http.StatusOK, "Orders fetched successfully", orders, total, int64(page), int64(limit))
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.Preload("User").Preload("Items.Product").Where("id = ?", id).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Order not found")
		return
	}

	// RBAC Check for Waiter
	role, _ := c.Get("role")
	if role == "waiter" {
		userID, _ := c.Get("user_id")
		if order.UserID != userID {
			utils.Error(c, http.StatusForbidden, "You can only view your own orders")
			return
		}
	}

	utils.Success(c, http.StatusOK, "Order fetched successfully", order)
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.Preload("Items").Where("id = ?", id).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Order not found")
		return
	}

	// RBAC Check for Waiter
	role, _ := c.Get("role")
	if role == "waiter" {
		userID, _ := c.Get("user_id")
		if order.UserID != userID {
			utils.Error(c, http.StatusForbidden, "You can only update your own orders")
			return
		}
	}

	var input UpdateOrderStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
	}

	if !validStatuses[input.Status] {
		utils.Error(c, http.StatusBadRequest, "Invalid status. Allowed statuses are: pending, completed")
		return
	}

	order.Status = input.Status
	if err := config.DB.Save(&order).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not update order status")
		return
	}

	utils.Success(c, http.StatusOK, "Order status updated successfully", order)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.Preload("Items").Where("id = ?", id).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Order not found")
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Restore stock if the order was pending (meaning stock was deducted but not yet consumed/completed)
		if order.Status == "pending" {
			for _, item := range order.Items {
				if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Where("order_id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Could not delete order and restore stock")
		return
	}

	utils.Success(c, http.StatusOK, "Order deleted successfully", nil)
}
