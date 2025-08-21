# Wildberries Order Service

Демонстрационный сервис для работы с заказами Wildberries с использованием Kafka, PostgreSQL и кеширования.

## 🚀 Быстрый старт

### Способ 1: Локальный запуск (без Docker)

```bash
# Перейти в директорию проекта
cd /workspace

# Установить зависимости
go mod tidy

# Запустить сервер (будет использован memory repository)
go run cmd/server/main.go
```

### Способ 2: С PostgreSQL через Docker

```bash
# Запустить PostgreSQL
docker compose up -d postgres

# В другом терминале запустить приложение
go run cmd/server/main.go
```

### Способ 3: Полная контейнеризация

```bash
# Запустить все сервисы в Docker
docker compose --profile full up -d
```

Сервер будет доступен по адресу: http://localhost:8081

📖 **Подробная документация**:
- [DOCKER.md](DOCKER.md) - Docker и контейнеризация
- [KAFKA.md](KAFKA.md) - Kafka интеграция и event-driven архитектура

### 📱 Интерфейсы

- **Веб-интерфейс**: http://localhost:8081
- **API для получения заказа**: `GET /order/{order_id}`
- **Health check**: `GET /health`

### 🧪 Тестирование

Для тестирования используйте тестовый ID заказа: `b563feb7b2b84b6test`

**Примеры запросов:**

```bash
# Через curl
curl http://localhost:8081/order/b563feb7b2b84b6test

# Через веб-интерфейс
# Откройте http://localhost:8081 и введите: b563feb7b2b84b6test
```

## 📁 Структура проекта

```
wildberries-order-service/
├── cmd/
│   ├── server/              # Основное приложение
│   │   └── main.go
│   └── producer/            # Kafka producer утилита
│       └── main.go
├── internal/
│   ├── config/              # Конфигурация приложения
│   │   └── config.go
│   ├── handlers/            # HTTP обработчики
│   │   ├── order.go         # API для заказов
│   │   └── web.go           # Веб-интерфейс
│   ├── kafka/               # Kafka интеграция
│   │   ├── consumer.go      # Kafka consumer
│   │   └── producer.go      # Kafka producer
│   ├── models/              # Модели данных и интерфейсы
│   │   └── order.go
│   ├── repository/          # Слой доступа к данным
│   │   ├── memory.go        # In-memory хранилище
│   │   └── postgres.go      # PostgreSQL хранилище
│   └── service/             # Бизнес-логика
│       └── order.go
├── migrations/              # SQL миграции
│   └── 001_create_tables.sql
├── docker-compose.yml       # Docker Compose конфигурация
├── Dockerfile              # Docker образ для приложения
├── .env                    # Переменные окружения
├── go.mod                  # Go модуль
├── README.md               # Основная документация
└── DOCKER.md               # Docker документация
```

## 🎯 Текущий статус

✅ **Реализовано:**
- Базовый HTTP сервер на порту 8081
- Модульная архитектура (models, repository, service, handlers)
- Эндпоинт для получения заказов по ID
- Эндпоинт для получения всех заказов
- Современный веб-интерфейс
- Хранилище заказов в памяти с thread-safety
- Health check эндпоинт
- Валидация данных заказов
- Graceful shutdown
- Dependency injection

✅ **Дополнительно реализовано:**
- PostgreSQL репозиторий с connection pooling
- Конфигурация через переменные окружения
- Docker и Docker Compose настройка
- Fallback на memory repository при недоступности БД
- Полная контейнеризация приложения
- pgAdmin для управления БД
- **Kafka Consumer** для получения сообщений о заказах
- **Kafka Producer** для тестирования и отправки сообщений
- **Event-driven архитектура** с автоматической обработкой
- **Graceful shutdown** для всех компонентов

🔄 **В разработке:**
- Улучшенное кеширование (Redis)
- Мониторинг и метрики (Prometheus, Grafana)
- Unit и integration тесты
- Kubernetes deployment манифесты

## 🧠 Что изучается в этом примере

### Архитектура Go приложений
- Модульная структура проекта
- Разделение ответственности (models, repository, service, handlers)
- Dependency Injection
- Интерфейсы для абстракции

### HTTP в Go
- Использование пакета `net/http`
- Создание HTTP handlers
- Работа с маршрутизацией
- Установка заголовков HTTP
- Обработка различных HTTP статусов

### JSON в Go
- Использование `encoding/json`
- Структуры с JSON тегами
- Сериализация и десериализация
- Обработка ошибок при работе с JSON

### Веб-разработка
- Создание современного HTML интерфейса
- Использование JavaScript для API вызовов
- Обработка пользовательского ввода
- Отображение результатов

### Concurrency и Thread Safety
- Использование sync.RWMutex для безопасного доступа к данным
- Graceful shutdown с использованием каналов и сигналов

### Валидация данных
- Создание валидаторов для бизнес-логики
- Обработка ошибок валидации
- Возврат понятных сообщений об ошибках

### Event-driven архитектура и Kafka
- Использование Apache Kafka для асинхронной обработки
- Consumer Groups и партиционирование
- Обработка ошибок и retry логика
- Мониторинг и отладка message flow
- Producer/Consumer patterns

## 📚 Полезные ресурсы

- [Go net/http документация](https://golang.org/pkg/net/http/)
- [JSON и Go](https://golang.org/blog/json)
- [HTTP Server примеры](https://gobyexample.com/http-servers)

## 🔄 Следующие шаги

1. Добавить подключение к PostgreSQL
2. Реализовать кеширование
3. Интегрировать Kafka
4. Добавить валидацию и обработку ошибок
5. Улучшить веб-интерфейс