package handler

import (
	"fmt"
	"net/http"

	"github.com/n-chetelat/garlic-service/services/common/genproto/orders"
	"github.com/n-chetelat/garlic-service/services/common/util"
	"github.com/n-chetelat/garlic-service/services/orders/types"
)

type OrdersHttpHandler struct {
	orderService types.OrderService
}

func NewHttpOrdersHandler(orderService types.OrderService) *OrdersHttpHandler {
	handler := &OrdersHttpHandler{
		orderService: orderService,
	}

	return handler
}

func (h *OrdersHttpHandler) RegisterRouter(router *http.ServeMux) {
	router.HandleFunc("POST /orders", h.CreateOrder)
	router.HandleFunc("DELETE /orders/", h.DeleteOrder)
}

func (h *OrdersHttpHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &orders.Order{
		CustomerId: req.GetCustomerId(),
		ProductId:  req.GetProductId(),
		Quantity:   req.GetQuantity(),
	}

	err = h.orderService.CreateOrder(r.Context(), order)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := &orders.CreateOrderResponse{Status: "success"}
	util.WriteJSON(w, http.StatusOK, res)
}

func (h *OrdersHttpHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("orderId is required"))
		return
	}

	err := h.orderService.DeleteOrder(r.Context(), orderId)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := &orders.DeleteOrderResponse{Status: "success"}
	util.WriteJSON(w, http.StatusOK, res)
}
