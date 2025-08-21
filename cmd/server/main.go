package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Временная структура для демонстрации
type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items []struct {
		ChrtID      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		RID         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmID        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}

// Временное хранилище заказов в памяти (позже заменим на кэш)
var orders = map[string]Order{
	"b563feb7b2b84b6test": {
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		}{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: struct {
			Transaction  string `json:"transaction"`
			RequestID    string `json:"request_id"`
			Currency     string `json:"currency"`
			Provider     string `json:"provider"`
			Amount       int    `json:"amount"`
			PaymentDt    int    `json:"payment_dt"`
			Bank         string `json:"bank"`
			DeliveryCost int    `json:"delivery_cost"`
			GoodsTotal   int    `json:"goods_total"`
			CustomFee    int    `json:"custom_fee"`
		}{
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
		Items: []struct {
			ChrtID      int    `json:"chrt_id"`
			TrackNumber string `json:"track_number"`
			Price       int    `json:"price"`
			RID         string `json:"rid"`
			Name        string `json:"name"`
			Sale        int    `json:"sale"`
			Size        string `json:"size"`
			TotalPrice  int    `json:"total_price"`
			NmID        int    `json:"nm_id"`
			Brand       string `json:"brand"`
			Status      int    `json:"status"`
		}{
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
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	},
}

// Обработчик для получения заказа по ID
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовки для JSON ответа
	w.Header().Set("Content-Type", "application/json")
	
	// Извлекаем ID заказа из URL
	path := strings.TrimPrefix(r.URL.Path, "/order/")
	orderID := path
	
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order ID is required",
		})
		return
	}
	
	// Ищем заказ в нашем временном хранилище
	order, exists := orders[orderID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order not found",
		})
		return
	}
	
	// Возвращаем заказ в JSON формате
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Error encoding order to JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
		return
	}
	
	log.Printf("Order %s retrieved successfully", orderID)
}

// Обработчик для главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Wildberries Order Service</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .search-form {
            margin: 20px 0;
            text-align: center;
        }
        input[type="text"] {
            padding: 10px;
            width: 300px;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-right: 10px;
        }
        button {
            padding: 10px 20px;
            background-color: #007bff;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
        button:hover {
            background-color: #0056b3;
        }
        .result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .result.success {
            background-color: #d4edda;
            border: 1px solid #c3e6cb;
        }
        .result.error {
            background-color: #f8d7da;
            border: 1px solid #f5c6cb;
        }
        pre {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🛒 Wildberries Order Service</h1>
        <p>Введите ID заказа для получения информации:</p>
        
        <div class="search-form">
            <input type="text" id="orderInput" placeholder="Введите Order ID (например: b563feb7b2b84b6test)" />
            <button onclick="searchOrder()">Найти заказ</button>
        </div>
        
        <div id="result" class="result"></div>
    </div>

    <script>
        function searchOrder() {
            const orderId = document.getElementById('orderInput').value.trim();
            const resultDiv = document.getElementById('result');
            
            if (!orderId) {
                showError('Пожалуйста, введите ID заказа');
                return;
            }
            
            // Показываем индикатор загрузки
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = '<p>Загрузка...</p>';
            
            // Выполняем запрос к API
            fetch('/order/' + encodeURIComponent(orderId))
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Заказ не найден');
                    }
                    return response.json();
                })
                .then(data => {
                    showSuccess(data);
                })
                .catch(error => {
                    showError('Ошибка: ' + error.message);
                });
        }
        
        function showSuccess(orderData) {
            const resultDiv = document.getElementById('result');
            resultDiv.className = 'result success';
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = 
                '<h3>✅ Заказ найден!</h3>' +
                '<pre>' + JSON.stringify(orderData, null, 2) + '</pre>';
        }
        
        function showError(message) {
            const resultDiv = document.getElementById('result');
            resultDiv.className = 'result error';
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = '<h3>❌ ' + message + '</h3>';
        }
        
        // Обработка Enter в поле ввода
        document.getElementById('orderInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchOrder();
            }
        });
    </script>
</body>
</html>`
	
	fmt.Fprint(w, html)
}

// Обработчик для проверки здоровья сервиса
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "wildberries-order-service",
	})
}

func main() {
	// Настраиваем маршруты
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/order/", getOrderHandler)
	http.HandleFunc("/health", healthHandler)
	
	port := "8081"
	log.Printf("🚀 Сервер запускается на порту %s", port)
	log.Printf("📱 Веб-интерфейс: http://localhost:%s", port)
	log.Printf("🔍 API endpoint: http://localhost:%s/order/{order_id}", port)
	log.Printf("❤️  Health check: http://localhost:%s/health", port)
	log.Printf("📝 Тестовый заказ: b563feb7b2b84b6test")
	
	// Запускаем сервер
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}