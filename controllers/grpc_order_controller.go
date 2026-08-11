package controllers

import (
	"context"
	"log"
	"restaurant-order-management/config"
	"restaurant-order-management/models"
	"restaurant-order-management/pb"
	"time"
)

type GrpcOrderServer struct {
	pb.UnimplementedOrderServiceServer
	UserClient pb.UserServiceClient
}

func (s *GrpcOrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	var order models.Order

	if err := config.DB.Where("id = ?", req.Id).First(&order).Error; err != nil {
		log.Println("Error fetching order, err: ", err)
		return nil, err
	}

	userRequest := pb.GetUserRequest{Id: order.UserID}
	user, err := s.UserClient.GetUser(ctx, &userRequest)
	if err != nil {
		log.Println("Cannot fetch user over gRPC:", err)
		return nil, err
	}

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

func (s *GrpcOrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {

	id := req.Id
	var order models.Order

	if err := config.DB.Preload("Items").Where("id = ?", id).First(&order).Error; err != nil {
		log.Println("Order not found")
		return nil, err
	}

	order.Status = req.Status

	if err := config.DB.Where("id = ?", id).Save(&order).Error; err != nil {
		log.Println("Order not updated")
		return nil, err
	}

	response := &pb.Order{
		Id:         order.Id,
		Status:     order.Status,
		Subtotal:   order.Subtotal,
		Tax:        order.Tax,
		GrandTotal: order.GrandTotal,
		UserId:     order.UserID,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  order.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil

}
