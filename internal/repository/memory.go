package repository

import (
	"errors"
	"sync"
	"time"

	"wildberries-order-service/internal/models"
)

// MemoryOrderRepository реализует хранение заказов в памяти
type MemoryOrderRepository struct {
	orders map[string]*models.Order
	mutex  sync.RWMutex
}

// NewMemoryOrderRepository создает новый репозиторий в памяти
func NewMemoryOrderRepository() *MemoryOrderRepository {
	repo := &MemoryOrderRepository{
		orders: make(map[string]*models.Order),
	}
	
	// Добавляем тестовые данные
	repo.initTestData()
	
	return repo
}

// GetByID возвращает заказ по ID
func (r *MemoryOrderRepository) GetByID(orderID string) (*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	order, exists := r.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	
	return order, nil
}

// Save сохраняет заказ
func (r *MemoryOrderRepository) Save(order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}
	
	if order.OrderUID == "" {
		return errors.New("order_uid is required")
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.orders[order.OrderUID] = order
	return nil
}

// GetAll возвращает все заказы
func (r *MemoryOrderRepository) GetAll() ([]*models.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	orders := make([]*models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	
	return orders, nil
}

// initTestData инициализирует тестовые данные
func (r *MemoryOrderRepository) initTestData() {
	testOrder := &models.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          "1",
	}
	
	r.orders[testOrder.OrderUID] = testOrder
}