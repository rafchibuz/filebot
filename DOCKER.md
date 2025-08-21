# 🐳 Docker Setup для Wildberries Order Service

Этот документ содержит инструкции по запуску сервиса с использованием Docker и Docker Compose.

## 📋 Предварительные требования

- Docker 20.10+
- Docker Compose 2.0+

### Установка Docker (если не установлен)

**Ubuntu/Debian:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**macOS:**
```bash
brew install --cask docker
```

**Windows:**
Скачайте Docker Desktop с официального сайта.

## 🚀 Быстрый старт

### 1. Запуск только PostgreSQL
```bash
# Запуск PostgreSQL и pgAdmin
docker compose up -d postgres pgadmin

# Проверка статуса
docker compose ps
```

### 2. Запуск всего стека
```bash
# Запуск PostgreSQL, pgAdmin и приложения
docker compose --profile full up -d

# Просмотр логов
docker compose logs -f app
```

### 3. Запуск для разработки
```bash
# Запуск только БД, приложение запускаем локально
docker compose up -d postgres

# В другом терминале
go run cmd/server/main.go
```

## 🔧 Конфигурация

### Переменные окружения

Создайте файл `.env` с настройками:

```env
# HTTP Server
PORT=8081

# PostgreSQL (для локального запуска)
DB_HOST=localhost
DB_PORT=5432
DB_USER=wb_user
DB_PASSWORD=wb_password
DB_NAME=wildberries_orders
DB_SSLMODE=disable

# PostgreSQL (для Docker контейнера)
# DB_HOST=postgres
```

### Файл docker-compose.yml

Основные сервисы:
- **postgres**: PostgreSQL 15 база данных
- **pgadmin**: Веб-интерфейс для управления БД
- **app**: Go приложение (запускается с профилем `full`)

## 📊 Управление базой данных

### pgAdmin
- URL: http://localhost:8080
- Email: admin@wildberries.local
- Password: admin123

### Подключение к PostgreSQL из pgAdmin
- Host: postgres (внутри Docker сети) или localhost (снаружи)
- Port: 5432
- Database: wildberries_orders
- Username: wb_user
- Password: wb_password

### Прямое подключение к PostgreSQL
```bash
# Подключение через psql
docker compose exec postgres psql -U wb_user -d wildberries_orders

# Или с хоста (если установлен postgresql-client)
psql -h localhost -p 5432 -U wb_user -d wildberries_orders
```

## 🛠️ Полезные команды

### Управление контейнерами
```bash
# Просмотр статуса
docker compose ps

# Просмотр логов
docker compose logs -f [service_name]

# Остановка сервисов
docker compose down

# Остановка с удалением volumes (ОСТОРОЖНО!)
docker compose down -v

# Пересборка образов
docker compose build --no-cache

# Запуск конкретного сервиса
docker compose up -d postgres
```

### Работа с данными
```bash
# Бэкап базы данных
docker compose exec postgres pg_dump -U wb_user wildberries_orders > backup.sql

# Восстановление базы данных
docker compose exec -T postgres psql -U wb_user wildberries_orders < backup.sql

# Очистка данных
docker compose exec postgres psql -U wb_user -d wildberries_orders -c "TRUNCATE orders CASCADE;"
```

### Мониторинг
```bash
# Использование ресурсов
docker compose top

# Статистика контейнеров
docker stats

# Проверка health check
docker compose exec app wget -qO- http://localhost:8081/health
```

## 🔍 Отладка

### Проблемы с подключением к БД
```bash
# Проверка сети
docker network ls
docker network inspect wildberries-order-service_wb_network

# Проверка подключения
docker compose exec app ping postgres
```

### Просмотр логов
```bash
# Все сервисы
docker compose logs

# Конкретный сервис
docker compose logs postgres
docker compose logs app

# Следить за логами в реальном времени
docker compose logs -f app
```

### Вход в контейнер
```bash
# PostgreSQL
docker compose exec postgres bash

# Приложение
docker compose exec app sh
```

## 🏗️ Сборка собственного образа

```bash
# Сборка образа
docker build -t wildberries-order-service:latest .

# Запуск собственного образа
docker run -p 8081:8081 --env-file .env wildberries-order-service:latest
```

## 🌐 Доступные эндпоинты

После запуска доступны:

- **Приложение**: http://localhost:8081
- **API**: http://localhost:8081/order/{id}
- **Health Check**: http://localhost:8081/health
- **pgAdmin**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## 🔒 Безопасность

В production рекомендуется:

1. Изменить пароли по умолчанию
2. Использовать Docker secrets
3. Настроить SSL/TLS
4. Ограничить сетевой доступ
5. Регулярно обновлять образы

```bash
# Пример с Docker secrets
echo "secure_password" | docker secret create db_password -
```

## 📈 Масштабирование

```bash
# Запуск нескольких экземпляров приложения
docker compose up -d --scale app=3

# С load balancer (требует настройки nginx/traefik)
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```