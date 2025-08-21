# 📨 Kafka Integration для Wildberries Order Service

Этот документ описывает интеграцию с Apache Kafka для обработки заказов в режиме реального времени.

## 🏗️ Архитектура

```
Kafka Producer → Kafka Topic (orders) → Kafka Consumer → Order Service → Database/Cache
```

**Компоненты:**
- **Kafka Consumer**: Получает сообщения из топика `orders`
- **Kafka Producer**: Отправляет тестовые сообщения (утилита)
- **Order Service**: Обрабатывает и валидирует заказы
- **Repository**: Сохраняет в PostgreSQL или Memory

## 🔧 Конфигурация

### Переменные окружения

```env
# Kafka Configuration
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=orders
KAFKA_GROUP_ID=wildberries-order-service
```

### Docker Compose

Kafka стек включает:
- **Zookeeper**: Координация Kafka кластера
- **Kafka**: Брокер сообщений
- **Kafka UI**: Веб-интерфейс для управления

```bash
# Запуск Kafka стека
docker compose up -d zookeeper kafka kafka-ui

# Проверка состояния
docker compose ps
```

## 📤 Kafka Producer

### Использование утилиты

```bash
# Отправить один тестовый заказ
go run cmd/producer/main.go

# Отправить несколько тестовых заказов
go run cmd/producer/main.go -test -count=5

# Отправить заказ из файла (планируется)
go run cmd/producer/main.go -file=order.json
```

### Программное использование

```go
// Создание producer
producer, err := kafka.NewProducer(&cfg.Kafka)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

// Отправка заказа
err = producer.SendOrder(order)
if err != nil {
    log.Printf("Failed to send order: %v", err)
}

// Асинхронная отправка
producer.SendOrderAsync(order)
```

## 📨 Kafka Consumer

### Автоматический запуск

Consumer автоматически запускается при старте приложения:

```go
// В main.go
kafkaConsumer, err := kafka.NewConsumer(&cfg.Kafka, orderService)
if err != nil {
    log.Printf("Failed to create Kafka consumer: %v", err)
} else if kafkaConsumer != nil {
    kafkaConsumer.Start()
    defer kafkaConsumer.Stop()
}
```

### Обработка сообщений

1. **Получение сообщения** из топика
2. **Десериализация** JSON в структуру Order
3. **Валидация** данных заказа
4. **Сохранение** через OrderService
5. **Подтверждение** обработки (commit offset)

### Обработка ошибок

- **Невалидные сообщения**: Логируются и пропускаются
- **Ошибки БД**: Сообщение не подтверждается (retry)
- **Ошибки парсинга**: Сообщение помечается как обработанное

## 🌐 Kafka UI

После запуска доступен веб-интерфейс:
- **URL**: http://localhost:8090
- **Функции**: 
  - Просмотр топиков
  - Мониторинг consumer groups
  - Отправка тестовых сообщений
  - Просмотр offset'ов

## 📊 Мониторинг

### Логи Consumer

```bash
# Логи приложения
docker compose logs -f app

# Логи Kafka
docker compose logs -f kafka
```

### Метрики

Consumer предоставляет информацию о:
- Количестве обработанных сообщений
- Ошибках обработки
- Статусе подключения

## 🧪 Тестирование

### 1. Локальное тестирование (без Kafka)

```bash
# Запуск с отключенным Kafka
KAFKA_ENABLED=false go run cmd/server/main.go
```

### 2. Тестирование с Kafka

```bash
# 1. Запуск Kafka
docker compose up -d zookeeper kafka kafka-ui

# 2. Запуск приложения
go run cmd/server/main.go

# 3. Отправка тестовых сообщений
go run cmd/producer/main.go -test -count=3
```

### 3. Проверка результатов

```bash
# Проверка API
curl http://localhost:8081/orders

# Проверка конкретного заказа
curl http://localhost:8081/order/{order_id}
```

## 🔄 Жизненный цикл сообщения

```
1. Producer отправляет Order JSON → Kafka Topic
2. Consumer получает сообщение
3. Парсинг JSON → Order struct
4. Валидация Order
5. OrderService.CreateOrder()
6. Сохранение в PostgreSQL/Memory
7. Commit offset (сообщение обработано)
```

## ⚙️ Настройки производительности

### Consumer

```go
config.Consumer.Group.Session.Timeout = 10 * time.Second
config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
config.Consumer.MaxProcessingTime = 2 * time.Minute
config.Consumer.Offsets.AutoCommit.Enable = true
```

### Producer

```go
config.Producer.RequiredAcks = sarama.WaitForAll
config.Producer.Retry.Max = 3
config.Producer.Compression = sarama.CompressionSnappy
config.Producer.Flush.Frequency = 500 * time.Millisecond
```

## 🚨 Обработка ошибок

### Consumer Errors

- **Connection Lost**: Автоматический переподключение
- **Invalid JSON**: Логирование + skip сообщения
- **Validation Failed**: Логирование + skip сообщения
- **Database Error**: Retry (не commit offset)

### Producer Errors

- **Connection Failed**: Ошибка отправки
- **Topic Not Found**: Автосоздание (если включено)
- **Serialization Error**: Ошибка до отправки

## 📈 Масштабирование

### Horizontal Scaling

```bash
# Запуск нескольких consumer instances
docker compose up -d --scale app=3
```

### Partitioning

```bash
# Создание топика с несколькими партициями
kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic orders \
  --partitions 3 \
  --replication-factor 1
```

## 🔍 Отладка

### Проверка топиков

```bash
# Список топиков
docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Описание топика
docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic orders
```

### Просмотр сообщений

```bash
# Consumer from beginning
docker compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic orders \
  --from-beginning
```

### Consumer Groups

```bash
# Список consumer groups
docker compose exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list

# Статус группы
docker compose exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group wildberries-order-service
```

## 🛡️ Безопасность

В production рекомендуется:

1. **SSL/TLS** шифрование
2. **SASL** аутентификация
3. **ACL** для контроля доступа
4. **Network policies** в Kubernetes

```yaml
# Пример SASL конфигурации
KAFKA_SASL_ENABLED_MECHANISMS: PLAIN
KAFKA_SASL_MECHANISM_INTER_BROKER_PROTOCOL: PLAIN
```

## 📚 Полезные команды

```bash
# Полная очистка Kafka данных
docker compose down -v
docker compose up -d zookeeper kafka

# Мониторинг потребления
watch -n 1 'docker compose exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group wildberries-order-service'

# Отправка сообщения через консоль
docker compose exec kafka kafka-console-producer --bootstrap-server localhost:9092 --topic orders
```