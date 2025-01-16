package handler

import (
	"context"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
	"github.com/n-chetelat/garlic-service/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	ordersService types.OrderService
	orders.UnimplementedOrderServiceServer
}

func NewGrpcOrdersHandler(grpc *grpc.Server, ordersService types.OrderService) {
	gRPCHandler := &OrdersGrpcHandler{
		ordersService: ordersService,
	}

	orders.RegisterOrderServiceServer(grpc, gRPCHandler)
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		CustomerId: req.GetCustomerId(),
		ProductId:  req.GetProductId(),
		Quantity:   req.GetQuantity(),
	}

	err := h.ordersService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	res := &orders.CreateOrderResponse{Status: "success"}

	return res, nil
}

func (h *OrdersGrpcHandler) GetOrders(ctx context.Context, req *orders.GetOrdersRequest) (*orders.GetOrderResponse, error) {
	orderCollection := h.ordersService.GetOrders(ctx)
	res := &orders.GetOrderResponse{Orders: orderCollection}
	return res, nil
}

func (h *OrdersGrpcHandler) DeleteOrder(ctx context.Context, req *orders.DeleteOrderRequest) (*orders.DeleteOrderResponse, error) {
	err := h.ordersService.DeleteOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	res := &orders.DeleteOrderResponse{Status: "success"}
	return res, nil
}
