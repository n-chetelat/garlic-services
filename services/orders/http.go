package main

import (
	"log"
	"net/http"

	handler "github.com/n-chetelat/garlic-service/services/orders/handler/orders"
	"github.com/n-chetelat/garlic-service/services/orders/service"
)

type httpServer struct {
	address string
}

func NewHttpServer(address string) *httpServer {
	return &httpServer{address: address}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	orderService := service.NewOrdersService()
	orderHandler := handler.NewHttpOrdersHandler(orderService)
	orderHandler.RegisterRouter(router)

	log.Println("Starting http server on ", s.address)

	return http.ListenAndServe(s.address, router)
}
