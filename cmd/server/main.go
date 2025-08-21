package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"wildberries-order-service/internal/config"
	"wildberries-order-service/internal/handlers"
	"wildberries-order-service/internal/kafka"
	"wildberries-order-service/internal/models"
	"wildberries-order-service/internal/repository"
	"wildberries-order-service/internal/service"
)

func main() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	
	// Загружаем конфигурацию
	cfg := config.Load()
	
	// Инициализация репозитория
	var orderRepo models.OrderRepository
	var err error
	
	// Пробуем подключиться к PostgreSQL, если не получается - используем memory
	postgresRepo, err := repository.NewPostgresOrderRepository(cfg.Database.GetConnectionString())
	if err != nil {
		log.Printf("⚠️  Failed to connect to PostgreSQL: %v", err)
		log.Println("🔄 Falling back to memory repository")
		orderRepo = repository.NewMemoryOrderRepository()
	} else {
		log.Println("🐘 Using PostgreSQL repository")
		orderRepo = postgresRepo
		
		// Настраиваем graceful shutdown для PostgreSQL
		defer func() {
			if pgRepo, ok := postgresRepo.(*repository.PostgresOrderRepository); ok {
				pgRepo.Close()
			}
		}()
	}
	
	// Инициализация сервиса
	orderService := service.NewOrderService(orderRepo)
	
	// Инициализация Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(&cfg.Kafka, orderService)
	if err != nil {
		log.Printf("⚠️  Failed to create Kafka consumer: %v", err)
		log.Println("🔄 Continuing without Kafka consumer")
	} else if kafkaConsumer != nil {
		// Запускаем consumer
		if err := kafkaConsumer.Start(); err != nil {
			log.Printf("⚠️  Failed to start Kafka consumer: %v", err)
		}
		
		// Настраиваем graceful shutdown для Kafka
		defer kafkaConsumer.Stop()
	}
	
	// Создание обработчиков
	orderHandler := handlers.NewOrderHandler(orderService)
	webHandler := handlers.NewWebHandler()
	
	// Настройка маршрутов
	setupRoutes(orderHandler, webHandler)
	
	// Запуск сервера
	port := cfg.Server.Port
	log.Printf("🚀 Сервер запускается на порту %s", port)
	log.Printf("📱 Веб-интерфейс: http://localhost:%s", port)
	log.Printf("🔍 API endpoint: http://localhost:%s/order/{order_id}", port)
	log.Printf("📋 Все заказы: http://localhost:%s/orders", port)
	log.Printf("❤️  Health check: http://localhost:%s/health", port)
	log.Printf("📝 Тестовый заказ: b563feb7b2b84b6test")
	log.Printf("🗄️  Database: %s:%d", cfg.Database.Host, cfg.Database.Port)
	
	// Graceful shutdown
	setupGracefulShutdown()
	
	// Запуск HTTP сервера
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// setupRoutes настраивает маршруты HTTP сервера
func setupRoutes(orderHandler *handlers.OrderHandler, webHandler *handlers.WebHandler) {
	// Создаем собственный маршрутизатор
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		
		switch r.URL.Path {
		case "/":
			log.Printf("Serving home page")
			webHandler.Home(w, r)
		case "/health":
			log.Printf("Serving health check")
			webHandler.Health(w, r)
		case "/orders":
			log.Printf("Serving orders list")
			orderHandler.GetAllOrders(w, r)
		default:
			// Проверяем, начинается ли путь с /order/
			if len(r.URL.Path) > 7 && r.URL.Path[:7] == "/order/" {
				log.Printf("Serving single order")
				orderHandler.GetOrder(w, r)
			} else if r.URL.Path == "/order" {
				if r.Method == http.MethodPost {
					log.Printf("Creating order")
					orderHandler.CreateOrder(w, r)
				} else {
					log.Printf("Invalid order request")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"Order ID is required. Use /order/{id} for GET requests"}`))
				}
			} else {
				// 404 для неизвестных путей
				log.Printf("Unknown path: %s", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"Endpoint not found"}`))
			}
		}
	})
}

// getPort возвращает порт для сервера
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return port
}

// setupGracefulShutdown настраивает graceful shutdown
func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		log.Println("\n🛑 Получен сигнал завершения. Завершение работы...")
		// Здесь можно добавить логику для корректного завершения работы
		// (закрытие соединений с БД, очистка ресурсов и т.д.)
		os.Exit(0)
	}()
}