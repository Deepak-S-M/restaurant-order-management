package handler

import (
	"fmt"
	"net/http"
	"restaurant-order-management/pb"
	"restaurant-order-management/utils"

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

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login User
//
// @Summary Login User
// @Description Login User
// @Tags Users
// @Accept json
// @Produce json
// @Param login	body	LoginInput	true	"login of user"
// @Success 200 {object} any
// @Router /login [post]
func (userHandler UserHandler) Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := userHandler.userServiceClient.Login(c.Request.Context(), &pb.GetUserLoginRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to Login user",
			"error":   err.Error(),
		})
		return
	}

	fmt.Println("handler Error: ", err)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login User API",
		"data":    user,
	})
}

// Getting User by Id
//
// @Summary Getting User by Id
// @Description Getting User by Id
// @Tags Users
// @Accept json
// @Produce json
// @Param id	path	string	true	"id of user"
// @Security BearerAuth
// @Success 200 {object} any
// @Router /api/users/{id} [get]
func (userHandler UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := userHandler.userServiceClient.GetUser(c.Request.Context(), &pb.GetUserRequest{Id: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
			"error":   err.Error(),
		})
		return
	}

	fmt.Println("handler Error: ", err)

	c.JSON(http.StatusOK, gin.H{
		"message": "Get User API",
		"data":    user,
	})
}
