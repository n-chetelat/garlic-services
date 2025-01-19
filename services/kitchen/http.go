package main

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
	"github.com/n-chetelat/garlic-service/services/common/util"
)

type httpServer struct {
	address string
	client  orders.OrderServiceClient
}

func NewHttpServer(address string) *httpServer {
	conn := NewGrpcClient(":8080")
	client := orders.NewOrderServiceClient(conn)
	return &httpServer{address: address, client: client}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("/", s.handleCreateOrder)

	log.Println("Starting http server on ", s.address)
	return http.ListenAndServe(s.address, router)
}

func (s *httpServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = s.client.CreateOrder(ctx, &req)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}

	res, err := s.client.GetOrders(ctx, &orders.GetOrdersRequest{})
	if err != nil {
		log.Fatalf("client error: %v", err)
	}

	templ := template.Must(template.New("orders").Parse(ordersTemplate))

	if err := templ.Execute(w, res.GetOrders()); err != nil {
		log.Fatalf("template error: %v", err)
	}
}

var ordersTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kitchen Orders</title>
</head>
<body>
    <h1>Orders List</h1>
    <table border="1">
        <tr>
            <th>Order ID</th>
            <th>Customer ID</th>
            <th>Quantity</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.OrderId}}</td>
            <td>{{.CustomerId}}</td>
            <td>{{.Quantity}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>`
