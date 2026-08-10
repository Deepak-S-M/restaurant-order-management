package main

import (
	"context"
	"fmt"
	"log"
	"restaurant-order-management/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getUser() {
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting server")
	}

	client := pb.NewUserServiceClient(conn)

	userRequest := pb.GetUserRequest{
		Id: "505cb159-cbf6-4152-9495-7fb3d7c4b4a6",
	}

	ctx := context.Background()

	data, err := client.GetUser(ctx, &userRequest)
	if err != nil {
		log.Println("Error in get order", err)
		return
	}
	fmt.Println(data)
}

func getOrder() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting server", err)
	}

	client := pb.NewOrderServiceClient(conn)

	orderRequest := pb.GetOrderRequest{
		Id: "ed37f9bc-b4d6-497a-bb7f-ea27d0ad54fb",
	}

	ctx := context.Background()

	data, err := client.GetOrder(ctx, &orderRequest)
	if err != nil {
		log.Fatalf("Cannot Fetch order. err: ", err)
	}

	fmt.Println(data)

}

func main() {
	fmt.Println("Client...")

	fmt.Println("-- User")
	getUser()

	fmt.Println("-- Order")
	getOrder()

}
