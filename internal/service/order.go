package service

import (
	"errors"
	"log"
	"strings"

	"wildberries-order-service/internal/models"
)

// OrderService реализует бизнес-логику для работы с заказами
type OrderService struct {
	repo models.OrderRepository
}

// NewOrderService создает новый сервис заказов
func NewOrderService(repo models.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

// GetOrder возвращает заказ по ID
func (s *OrderService) GetOrder(orderID string) (*models.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		log.Printf("Failed to get order %s: %v", orderID, err)
		return nil, err
	}
	
	log.Printf("Order %s retrieved successfully", orderID)
	return order, nil
}

// CreateOrder создает новый заказ
func (s *OrderService) CreateOrder(order *models.Order) error {
	if err := s.ValidateOrder(order); err != nil {
		return err
	}
	
	if err := s.repo.Save(order); err != nil {
		log.Printf("Failed to save order %s: %v", order.OrderUID, err)
		return err
	}
	
	log.Printf("Order %s created successfully", order.OrderUID)
	return nil
}

// ValidateOrder валидирует заказ
func (s *OrderService) ValidateOrder(order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	
	if strings.TrimSpace(order.OrderUID) == "" {
		return errors.New("order_uid is required")
	}
	
	if strings.TrimSpace(order.TrackNumber) == "" {
		return errors.New("track_number is required")
	}
	
	// Валидация доставки
	if err := s.validateDelivery(&order.Delivery); err != nil {
		return err
	}
	
	// Валидация платежа
	if err := s.validatePayment(&order.Payment); err != nil {
		return err
	}
	
	// Валидация товаров
	if len(order.Items) == 0 {
		return errors.New("order must contain at least one item")
	}
	
	for i, item := range order.Items {
		if err := s.validateItem(&item); err != nil {
			return errors.New("item " + string(rune(i)) + ": " + err.Error())
		}
	}
	
	return nil
}

// validateDelivery валидирует данные доставки
func (s *OrderService) validateDelivery(delivery *models.Delivery) error {
	if strings.TrimSpace(delivery.Name) == "" {
		return errors.New("delivery name is required")
	}
	
	if strings.TrimSpace(delivery.Phone) == "" {
		return errors.New("delivery phone is required")
	}
	
	if strings.TrimSpace(delivery.Address) == "" {
		return errors.New("delivery address is required")
	}
	
	return nil
}

// validatePayment валидирует данные платежа
func (s *OrderService) validatePayment(payment *models.Payment) error {
	if strings.TrimSpace(payment.Transaction) == "" {
		return errors.New("payment transaction is required")
	}
	
	if strings.TrimSpace(payment.Currency) == "" {
		return errors.New("payment currency is required")
	}
	
	if payment.Amount <= 0 {
		return errors.New("payment amount must be positive")
	}
	
	return nil
}

// validateItem валидирует товар
func (s *OrderService) validateItem(item *models.Item) error {
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("item name is required")
	}
	
	if item.Price <= 0 {
		return errors.New("item price must be positive")
	}
	
	return nil
}

// GetAllOrders возвращает все заказы (для административных целей)
func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	orders, err := s.repo.GetAll()
	if err != nil {
		log.Printf("Failed to get all orders: %v", err)
		return nil, err
	}
	
	log.Printf("Retrieved %d orders", len(orders))
	return orders, nil
}