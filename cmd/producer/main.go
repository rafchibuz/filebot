package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"wildberries-order-service/internal/config"
	"wildberries-order-service/internal/kafka"
	"wildberries-order-service/internal/models"
)

func main() {
	// Флаги командной строки
	var (
		testMode = flag.Bool("test", false, "Send test order")
		count    = flag.Int("count", 1, "Number of test orders to send")
		file     = flag.String("file", "", "JSON file with order data")
	)
	flag.Parse()

	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем Kafka producer
	producer, err := kafka.NewProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	if producer == nil {
		log.Fatal("Kafka producer is disabled")
	}
	defer producer.Close()

	log.Printf("🚀 Kafka Producer started")
	log.Printf("📤 Brokers: %v", cfg.Kafka.Brokers)
	log.Printf("📋 Topic: %s", cfg.Kafka.Topic)

	if *testMode {
		// Отправляем тестовые заказы
		log.Printf("📦 Sending %d test orders...", *count)
		for i := 0; i < *count; i++ {
			if err := producer.SendTestOrder(); err != nil {
				log.Printf("❌ Failed to send test order %d: %v", i+1, err)
			} else {
				log.Printf("✅ Test order %d sent successfully", i+1)
			}
			
			// Небольшая пауза между отправками
			if i < *count-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	} else if *file != "" {
		// Отправляем заказ из файла
		log.Printf("📄 Sending order from file: %s", *file)
		if err := sendOrderFromFile(producer, *file); err != nil {
			log.Fatalf("❌ Failed to send order from file: %v", err)
		}
		log.Println("✅ Order from file sent successfully")
	} else {
		// Отправляем предопределенный заказ
		log.Println("📦 Sending predefined test order...")
		order := createSampleOrder()
		if err := producer.SendOrder(order); err != nil {
			log.Fatalf("❌ Failed to send order: %v", err)
		}
		log.Println("✅ Predefined order sent successfully")
	}

	log.Println("🏁 Producer finished")
}

func sendOrderFromFile(producer *kafka.Producer, filename string) error {
	// Здесь можно добавить чтение из файла
	// Пока используем заглушку
	log.Printf("📄 File reading not implemented yet: %s", filename)
	return fmt.Errorf("file reading not implemented")
}

func createSampleOrder() *models.Order {
	now := time.Now()
	timestamp := now.Format("20060102_150405")
	
	return &models.Order{
		OrderUID:    "sample_order_" + timestamp,
		TrackNumber: "SAMPLE_TRACK_" + timestamp,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Иван Иванов",
			Phone:   "+7 900 123-45-67",
			Zip:     "123456",
			City:    "Москва",
			Address: "ул. Тестовая, д. 1, кв. 10",
			Region:  "Московская область",
			Email:   "ivan.ivanov@example.com",
		},
		Payment: models.Payment{
			Transaction:  "sample_txn_" + timestamp,
			RequestID:    "",
			Currency:     "RUB",
			Provider:     "wbpay",
			Amount:       2500,
			PaymentDt:    now.Unix(),
			Bank:         "sberbank",
			DeliveryCost: 500,
			GoodsTotal:   2000,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      1234567,
				TrackNumber: "SAMPLE_TRACK_" + timestamp,
				Price:       1000,
				RID:         "sample_rid_1_" + timestamp,
				Name:        "Тестовый товар 1",
				Sale:        10,
				Size:        "L",
				TotalPrice:  900,
				NmID:        9876543,
				Brand:       "TestBrand",
				Status:      202,
			},
			{
				ChrtID:      7654321,
				TrackNumber: "SAMPLE_TRACK_" + timestamp,
				Price:       1200,
				RID:         "sample_rid_2_" + timestamp,
				Name:        "Тестовый товар 2",
				Sale:        0,
				Size:        "M",
				TotalPrice:  1100,
				NmID:        3456789,
				Brand:       "AnotherBrand",
				Status:      202,
			},
		},
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        "sample_customer_" + timestamp,
		DeliveryService:   "cdek",
		Shardkey:          "1",
		SmID:              99,
		DateCreated:       now,
		OofShard:          "1",
	}
}

func printOrderJSON(order *models.Order) {
	orderJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal order: %v", err)
		return
	}
	
	fmt.Println("📋 Order JSON:")
	fmt.Println(string(orderJSON))
}