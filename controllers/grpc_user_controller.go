package controllers

import (
	"context"
	"fmt"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/pb"
	"restaurant-order-management/utils"
)

type GrpcUserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *GrpcUserServer) Login(ctx context.Context, req *pb.GetUserLoginRequest) (*pb.UserLoginResponse, error) {
	var user models.User
	if err := config.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		fmt.Println("Invalid email or password")
		return nil, err
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		fmt.Println("Invalid email or password")
		return nil, err
	}

	token, err := utils.GenerateToken(user.Id, user.Email, user.Role.Name)
	if err != nil {
		fmt.Println("Could not generate token")
		return nil, err
	}

	response := &pb.UserLoginResponse{
		Token: token,
		User: &pb.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Name,
		},
	}

	return response, nil
}

func (s *GrpcUserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	var user models.User
	if err := config.DB.Preload("Role").Where("id = ?", req.Id).First(&user).Error; err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	response := &pb.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role.Name,
	}

	return response, nil
}
