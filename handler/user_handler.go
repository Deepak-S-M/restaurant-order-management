package handler

import (
	"fmt"
	"net/http"
	"restaurant-order-management/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type UserHandler struct {
	userServiceClient pb.UserServiceClient
}

func NewUserServiceHandler(conn *grpc.ClientConn) UserHandler {
	userClient := pb.NewUserServiceClient(conn)
	return UserHandler{userServiceClient: userClient}
}

func (userHandler UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := userHandler.userServiceClient.GetUser(c.Request.Context(), &pb.GetUserRequest{Id: userID})
	if err != nil {
		return
	}

	fmt.Println("Error: ", err)

	c.JSON(http.StatusOK, gin.H{
		"message": "Get User API",
		"data":    user,
	})
}
