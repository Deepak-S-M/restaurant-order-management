package controllertest

import (
	"context"
	"testing"

	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	pb "restaurant-order-management/pb"
	"restaurant-order-management/testutil"

	"github.com/go-openapi/testify/v2/require"
)

func TestGetUser_Success(t *testing.T) {

	role := testutil.RandomRole()
	user := testutil.RandomUser()
	user.Role = role
	err := config.DB.Create(&user).Error
	require.Nil(t, err)

	server := &controllers.GrpcUserServer{}
	resp, err := server.GetUser(context.Background(), &pb.GetUserRequest{Id: user.Id})

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.Id, resp.Id)
	require.Equal(t, user.Email, resp.Email)
	require.Equal(t, user.Name, resp.Name)
	require.Equal(t, user.Role.Name, resp.Role)
}

func TestGetUser_NotFound(t *testing.T) {
	server := &controllers.GrpcUserServer{}
	resp, err := server.GetUser(context.Background(), &pb.GetUserRequest{Id: "nonexistent-id"})

	require.NotNil(t, err)
	require.Nil(t, resp)
}
