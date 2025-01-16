package main

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
)

type httpServer struct {
	address string
}

func NewHttpServer(address string) *httpServer {
	return &httpServer{address: address}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	conn := NewGrpcClient(":8080")
	defer conn.Close()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := orders.NewOrderServiceClient(conn)

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()
		order := &orders.CreateOrderRequest{
			CustomerId: 4,
			ProductId:  5,
			Quantity:   90,
		}
		_, err := client.CreateOrder(ctx, order)
		if err != nil {
			log.Fatalf("client error: %v", err)
		}

		res, err := client.GetOrders(ctx, &orders.GetOrdersRequest{
			CustomerID: 5,
		})
		if err != nil {
			log.Fatalf("client error: %v", err)
		}

		templ := template.Must(template.New("orders").Parse(ordersTemplate))

		if err := templ.Execute(w, res.GetOrders()); err != nil {
			log.Fatalf("template error: %v", err)
		}
	})

	log.Println("Starting http server on ", s.address)

	return http.ListenAndServe(s.address, router)
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
