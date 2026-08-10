package controllers

import (
	"context"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/pb"
)

type GrpcUserServer struct {
	pb.UnimplementedUserServiceServer
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
