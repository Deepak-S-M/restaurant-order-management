package controllers

import (
	"context"
	"fmt"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/pb"
	"time"
)

type GrpcOrderServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *GrpcOrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {

	var order models.Order

	if err := config.DB.Preload("User").Preload("User.Role").Where("id = ?", req.Id).First(&order).Error; err != nil {
		log.Println("Error fetching order, err: ", err)
		return nil, err
	}

	user := &pb.User{
		Id:    order.User.Id,
		Email: order.User.Email,
		Name:  order.User.Name,
		Role:  order.User.Role.Name,
	}
	fmt.Println("user fetch: ", user)

	response := &pb.Order{
		Id:         order.Id,
		Status:     order.Status,
		Subtotal:   order.Subtotal,
		Tax:        order.Tax,
		GrandTotal: order.GrandTotal,
		UserId:     order.UserID,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
		User:       user,
	}

	return response, nil
}
