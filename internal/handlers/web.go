package handlers

import (
	"fmt"
	"net/http"
)

// WebHandler обрабатывает веб-интерфейс
type WebHandler struct{}

// NewWebHandler создает новый обработчик веб-интерфейса
func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

// Home обрабатывает главную страницу
func (h *WebHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Wildberries Order Service</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #ff6b6b, #feca57);
            padding: 40px;
            text-align: center;
            color: white;
        }
        
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }
        
        .content {
            padding: 40px;
        }
        
        .search-section {
            margin-bottom: 30px;
        }
        
        .search-section h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 1.5rem;
        }
        
        .search-form {
            display: flex;
            gap: 15px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .search-input {
            flex: 1;
            min-width: 300px;
            padding: 15px 20px;
            border: 2px solid #e1e8ed;
            border-radius: 12px;
            font-size: 1rem;
            transition: all 0.3s ease;
        }
        
        .search-input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        .search-btn {
            padding: 15px 30px;
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            border: none;
            border-radius: 12px;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            white-space: nowrap;
        }
        
        .search-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.3);
        }
        
        .search-btn:active {
            transform: translateY(0);
        }
        
        .example-orders {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 12px;
            margin-bottom: 20px;
        }
        
        .example-orders h3 {
            color: #495057;
            margin-bottom: 15px;
            font-size: 1.2rem;
        }
        
        .example-id {
            background: white;
            padding: 12px 16px;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.9rem;
            color: #495057;
            border: 1px solid #dee2e6;
            cursor: pointer;
            transition: all 0.2s ease;
        }
        
        .example-id:hover {
            background: #e9ecef;
            border-color: #adb5bd;
        }
        
        .result {
            margin-top: 30px;
            padding: 0;
            border-radius: 12px;
            display: none;
            overflow: hidden;
        }
        
        .result.show {
            display: block;
            animation: slideIn 0.3s ease;
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        .result.success {
            background: #d1f2eb;
            border: 2px solid #52c41a;
        }
        
        .result.error {
            background: #ffeaa7;
            border: 2px solid #fdcb6e;
        }
        
        .result.loading {
            background: #e3f2fd;
            border: 2px solid #2196f3;
        }
        
        .result-header {
            padding: 20px;
            font-weight: 600;
            font-size: 1.1rem;
        }
        
        .result.success .result-header {
            background: #52c41a;
            color: white;
        }
        
        .result.error .result-header {
            background: #fdcb6e;
            color: #333;
        }
        
        .result.loading .result-header {
            background: #2196f3;
            color: white;
        }
        
        .result-content {
            padding: 20px;
        }
        
        .json-viewer {
            background: #282c34;
            color: #abb2bf;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.85rem;
            line-height: 1.5;
            overflow-x: auto;
            white-space: pre;
        }
        
        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #2196f3;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin-right: 10px;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .api-info {
            background: #e8f4f8;
            padding: 20px;
            border-radius: 12px;
            margin-top: 30px;
            border-left: 4px solid #17a2b8;
        }
        
        .api-info h3 {
            color: #17a2b8;
            margin-bottom: 15px;
        }
        
        .api-endpoint {
            background: white;
            padding: 10px 15px;
            border-radius: 6px;
            font-family: monospace;
            margin: 5px 0;
            border: 1px solid #bee5eb;
        }
        
        @media (max-width: 768px) {
            .search-form {
                flex-direction: column;
            }
            
            .search-input {
                min-width: auto;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            .content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛒 Wildberries Order Service</h1>
            <p>Демонстрационный сервис для работы с заказами</p>
        </div>
        
        <div class="content">
            <div class="search-section">
                <h2>🔍 Поиск заказа</h2>
                <div class="search-form">
                    <input 
                        type="text" 
                        id="orderInput" 
                        class="search-input"
                        placeholder="Введите Order ID для поиска..."
                        autocomplete="off"
                    />
                    <button onclick="searchOrder()" class="search-btn">
                        Найти заказ
                    </button>
                </div>
                
                <div class="example-orders">
                    <h3>📝 Тестовые данные:</h3>
                    <div class="example-id" onclick="fillOrderId('b563feb7b2b84b6test')">
                        b563feb7b2b84b6test
                    </div>
                </div>
            </div>
            
            <div id="result" class="result"></div>
            
            <div class="api-info">
                <h3>🔗 API Endpoints</h3>
                <div class="api-endpoint">GET /order/{order_id} - Получить заказ по ID</div>
                <div class="api-endpoint">GET /orders - Получить все заказы</div>
                <div class="api-endpoint">GET /health - Проверка состояния сервиса</div>
                <div class="api-endpoint">POST /order - Создать новый заказ</div>
            </div>
        </div>
    </div>

    <script>
        function fillOrderId(orderId) {
            document.getElementById('orderInput').value = orderId;
            document.getElementById('orderInput').focus();
        }
        
        function searchOrder() {
            const orderId = document.getElementById('orderInput').value.trim();
            const resultDiv = document.getElementById('result');
            
            if (!orderId) {
                showError('Пожалуйста, введите ID заказа');
                return;
            }
            
            // Показываем индикатор загрузки
            showLoading('Поиск заказа...');
            
            // Выполняем запрос к API
            fetch('/order/' + encodeURIComponent(orderId))
                .then(response => {
                    if (!response.ok) {
                        if (response.status === 404) {
                            throw new Error('Заказ не найден');
                        } else {
                            throw new Error('Ошибка сервера');
                        }
                    }
                    return response.json();
                })
                .then(data => {
                    showSuccess(data);
                })
                .catch(error => {
                    showError(error.message);
                });
        }
        
        function showLoading(message) {
            const resultDiv = document.getElementById('result');
            resultDiv.className = 'result loading show';
            resultDiv.innerHTML = 
                '<div class="result-header">' +
                    '<div class="loading-spinner"></div>' + message +
                '</div>';
        }
        
        function showSuccess(orderData) {
            const resultDiv = document.getElementById('result');
            resultDiv.className = 'result success show';
            
            const jsonString = JSON.stringify(orderData, null, 2);
            
            resultDiv.innerHTML = 
                '<div class="result-header">✅ Заказ найден!</div>' +
                '<div class="result-content">' +
                    '<div class="json-viewer">' + escapeHtml(jsonString) + '</div>' +
                '</div>';
        }
        
        function showError(message) {
            const resultDiv = document.getElementById('result');
            resultDiv.className = 'result error show';
            resultDiv.innerHTML = 
                '<div class="result-header">❌ ' + escapeHtml(message) + '</div>';
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        // Обработка Enter в поле ввода
        document.getElementById('orderInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchOrder();
            }
        });
        
        // Автофокус на поле ввода
        document.getElementById('orderInput').focus();
    </script>
</body>
</html>`
	
	fmt.Fprint(w, html)
}

// Health обрабатывает проверку состояния сервиса
func (h *WebHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	fmt.Fprintf(w, `{"status":"healthy","service":"wildberries-order-service","version":"1.0.0"}`)
}