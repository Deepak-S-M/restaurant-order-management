package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restaurant-order-management/pb"
	"restaurant-order-management/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	orderServiceClient pb.OrderServiceClient
}

func NewOrderServiceHandler(conn *grpc.ClientConn) OrderHandler {
	orderClient := pb.NewOrderServiceClient(conn)
	return OrderHandler{orderServiceClient: orderClient}
}

var (
	sseClients = make(map[string][]chan *pb.Order)
	sseMutex   sync.Mutex
)

func notifySSEClients(
	orderID string,
	order *pb.Order,
) {

	sseMutex.Lock()
	defer sseMutex.Unlock()

	for _, clientChan := range sseClients[orderID] {

		select {

		case clientChan <- order:

		default:
		}
	}
}

func (orderHandler OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "SSE not supported"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	order, err := orderHandler.orderServiceClient.GetOrder(c.Request.Context(), &pb.GetOrderRequest{Id: orderID})
	if err != nil {
		errorData := map[string]interface{}{
			"message": "Failed to get order",
			"error":   err.Error(),
		}

		data, _ := json.Marshal(errorData)

		fmt.Fprintf(c.Writer, "event: error\n")
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)

		flusher.Flush()

		return
	}

	eventChan := make(chan *pb.Order)

	sseMutex.Lock()
	sseClients[orderID] = append(
		sseClients[orderID],
		eventChan,
	)
	sseMutex.Unlock()

	defer func() {

		sseMutex.Lock()
		defer sseMutex.Unlock()

		clients := sseClients[orderID]

		for i, ch := range clients {
			if ch == eventChan {
				sseClients[orderID] = append(
					clients[:i],
					clients[i+1:]...,
				)
				break
			}
		}

		if len(sseClients[orderID]) == 0 {
			delete(sseClients, orderID)
		}

		close(eventChan)
	}()

	data, _ := json.Marshal(order)

	fmt.Fprintf(c.Writer, "event: order\n")
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)

	flusher.Flush()

	for {

		select {

		case updatedOrder := <-eventChan:

			data, _ := json.Marshal(updatedOrder)

			fmt.Fprintf(c.Writer, "event: order\n")
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)

			flusher.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func (orderHandler OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var input UpdateOrderStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	validStatuses := map[string]bool{
		"pending": true,
		"ready":   true,
	}

	if !validStatuses[input.Status] {
		utils.Error(c, http.StatusBadRequest, "Invalid status. Allowed statuses are: pending, completed")
		return
	}

	order, err := orderHandler.orderServiceClient.UpdateOrderStatus(c.Request.Context(), &pb.UpdateOrderStatusRequest{Id: orderID, Status: input.Status})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update order",
			"error":   err.Error(),
		})
		return
	}

	notifySSEClients(orderID, order)

	c.JSON(http.StatusOK, gin.H{
		"message": "Update Order Status",
		"data":    order,
	})
}
