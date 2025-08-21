package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"wildberries-order-service/internal/config"
	"wildberries-order-service/internal/models"
)

// Producer представляет Kafka producer для отправки заказов
type Producer struct {
	config   *config.KafkaConfig
	producer sarama.SyncProducer
}

// NewProducer создает новый Kafka producer
func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	if !cfg.Enabled {
		log.Println("📵 Kafka producer disabled")
		return nil, nil
	}

	// Конфигурация Sarama
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	// Создаем producer
	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	log.Printf("📤 Kafka producer created for brokers: %v", cfg.Brokers)

	return &Producer{
		config:   cfg,
		producer: producer,
	}, nil
}

// SendOrder отправляет заказ в Kafka
func (p *Producer) SendOrder(order *models.Order) error {
	if p == nil {
		log.Println("📵 Kafka producer is disabled, skipping send")
		return nil
	}

	// Сериализуем заказ в JSON
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Создаем сообщение
	message := &sarama.ProducerMessage{
		Topic: p.config.Topic,
		Key:   sarama.StringEncoder(order.OrderUID),
		Value: sarama.ByteEncoder(orderBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	// Отправляем сообщение
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		log.Printf("❌ Failed to send message: %v", err)
		return err
	}

	log.Printf("✅ Message sent successfully: topic=%s, partition=%d, offset=%d, key=%s", 
		p.config.Topic, partition, offset, order.OrderUID)

	return nil
}

// SendOrderAsync отправляет заказ асинхронно (для демонстрации)
func (p *Producer) SendOrderAsync(order *models.Order) {
	if p == nil {
		log.Println("📵 Kafka producer is disabled, skipping async send")
		return
	}

	go func() {
		if err := p.SendOrder(order); err != nil {
			log.Printf("❌ Failed to send order async: %v", err)
		}
	}()
}

// Close закрывает producer
func (p *Producer) Close() {
	if p == nil {
		return
	}

	log.Println("🛑 Closing Kafka producer...")
	
	if err := p.producer.Close(); err != nil {
		log.Printf("❌ Error closing producer: %v", err)
	}
	
	log.Println("✅ Kafka producer closed")
}

// SendTestOrder отправляет тестовый заказ
func (p *Producer) SendTestOrder() error {
	if p == nil {
		return nil
	}

	testOrder := &models.Order{
		OrderUID:    "test_" + time.Now().Format("20060102_150405"),
		TrackNumber: "TEST_TRACK_" + time.Now().Format("150405"),
		Entry:       "TEST",
		Delivery: models.Delivery{
			Name:    "Test Customer",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "Test City",
			Address: "123 Test Street",
			Region:  "Test Region",
			Email:   "test@example.com",
		},
		Payment: models.Payment{
			Transaction:  "test_transaction_" + time.Now().Format("150405"),
			RequestID:    "",
			Currency:     "USD",
			Provider:     "test_provider",
			Amount:       1000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      123456,
				TrackNumber: "TEST_TRACK_" + time.Now().Format("150405"),
				Price:       900,
				RID:         "test_rid_" + time.Now().Format("150405"),
				Name:        "Test Product",
				Sale:        0,
				Size:        "M",
				TotalPrice:  900,
				NmID:        789012,
				Brand:       "Test Brand",
				Status:      200,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "test_delivery",
		Shardkey:          "1",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	return p.SendOrder(testOrder)
}