package web

import (
	"encoding/json"
	"net/http"

	"github.com/devfullcycle/fc-clean-architecture/internal/usecase"
	"github.com/devfullcycle/fc-clean-architecture/pkg/events"
)

// OrderHandler handles HTTP requests for order operations.
type OrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   interface{ Save(order interface{}) error }
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

// NewOrderHandler creates an OrderHandler with the provided use cases.
func NewOrderHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
) *OrderHandler {
	return &OrderHandler{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

// Create handles POST /order — creates a new order.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.CreateOrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// List handles GET /order — returns all orders.
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
