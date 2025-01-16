package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
)

var (
	ordersDB   = make([]*orders.Order, 0)
	orderIdSeq = int32(1)
	mu         sync.Mutex
)

type OrderService struct {
}

func NewOrdersService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	mu.Lock()
	defer mu.Unlock()

	order.OrderId = orderIdSeq
	orderIdSeq++

	ordersDB = append(ordersDB, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context) []*orders.Order {
	return ordersDB
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderId int32) error {
	for i, order := range ordersDB {
		if order.OrderId == orderId {
			ordersDB = append(ordersDB[:i], ordersDB[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found")
}
