package main

import (
	"log"
	"net"

	handler "github.com/n-chetelat/garlic-service/services/orders/handler/orders"
	"github.com/n-chetelat/garlic-service/services/orders/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	address string
}

func NewGRPCServer(address string) *gRPCServer {
	return &(gRPCServer{address: address})
}

func (s *gRPCServer) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	gRPCServer := grpc.NewServer()

	// register grpc srevices
	orderService := service.NewOrdersService()
	handler.NewGrpcOrdersHandler(gRPCServer, orderService)

	log.Println("Starting gRPC server on ", s.address)

	return gRPCServer.Serve(listener)
}
