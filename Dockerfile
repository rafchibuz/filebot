# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod и sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Используем минимальный образ для production
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из builder стадии
COPY --from=builder /app/main .

# Копируем .env файл (опционально)
COPY --from=builder /app/.env .

# Меняем владельца файлов
RUN chown -R appuser:appgroup /root/

# Переключаемся на обычного пользователя
USER appuser

# Открываем порт
EXPOSE 8081

# Устанавливаем health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Запускаем приложение
CMD ["./main"]