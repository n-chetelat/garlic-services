package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
	"github.com/oklog/ulid/v2"
)

var (
	ordersDB = make([]*orders.Order, 0)
	mu       sync.Mutex
)

type OrderService struct {
}

func NewOrdersService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	mu.Lock()
	defer mu.Unlock()

	order.OrderId = ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()

	ordersDB = append(ordersDB, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context) []*orders.Order {
	return ordersDB
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *orders.Order) error {
	mu.Lock()
	defer mu.Unlock()

	for i, existingOrder := range ordersDB {
		if existingOrder.OrderId == order.OrderId {
			ordersDB[i] = order
			return nil
		}
	}
	return fmt.Errorf("order not found")
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderId string) error {
	for i, order := range ordersDB {
		if order.OrderId == orderId {
			ordersDB = append(ordersDB[:i], ordersDB[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found")
}
