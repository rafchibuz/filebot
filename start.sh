#!/bin/bash

echo "🚀 Запуск Twitch Clips Aggregator..."

# Проверяем наличие .env файла
if [ ! -f ".env" ]; then
    echo "❌ Файл .env не найден. Запустите setup.sh для первоначальной настройки."
    exit 1
fi

# Проверяем наличие Twitch API ключей
if ! grep -q "TWITCH_CLIENT_ID=your_client_id_here" .env; then
    echo "✅ Twitch API ключи настроены"
else
    echo "⚠️  Пожалуйста, настройте Twitch API ключи в .env файле!"
    echo "Откройте .env файл и замените your_client_id_here и your_client_secret_here на ваши реальные ключи."
    exit 1
fi

# Проверяем установку зависимостей
if [ ! -d "node_modules" ]; then
    echo "📦 Установка зависимостей сервера..."
    npm install
fi

if [ ! -d "client/node_modules" ]; then
    echo "📦 Установка зависимостей клиента..."
    cd client && npm install && cd ..
fi

# Создание папки для базы данных
mkdir -p data

echo "🎬 Запуск приложения..."
echo "🔗 Клиент: http://localhost:3000"
echo "🔗 API: http://localhost:3001"
echo ""
echo "Нажмите Ctrl+C для остановки"

# Запуск приложения
npm run dev