package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"wildberries-order-service/internal/models"
)

// OrderHandler обрабатывает HTTP запросы для заказов
type OrderHandler struct {
	service models.OrderService
}

// NewOrderHandler создает новый обработчик заказов
func NewOrderHandler(service models.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// GetOrder обрабатывает GET /order/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для JSON ответа
	w.Header().Set("Content-Type", "application/json")
	
	// Извлекаем ID заказа из URL
	path := strings.TrimPrefix(r.URL.Path, "/order/")
	orderID := strings.TrimSpace(path)
	
	if orderID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Order ID is required")
		return
	}
	
	// Получаем заказ через сервис
	order, err := h.service.GetOrder(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			h.writeErrorResponse(w, http.StatusNotFound, "Order not found")
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	
	// Возвращаем заказ в JSON формате
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// CreateOrder обрабатывает POST /order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	
	if err := h.service.CreateOrder(&order); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Order created successfully",
		"order_id": order.OrderUID,
	})
}

// GetAllOrders обрабатывает GET /orders (для административных целей)
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	orders, err := h.service.GetAllOrders()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// writeErrorResponse записывает ошибку в ответ
func (h *OrderHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}