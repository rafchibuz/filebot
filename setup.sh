#!/bin/bash

echo "🚀 Настройка Twitch Clips Aggregator..."

# Проверяем наличие Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js не установлен. Пожалуйста, установите Node.js версии 16 или выше."
    exit 1
fi

# Проверяем версию Node.js
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 16 ]; then
    echo "❌ Требуется Node.js версии 16 или выше. Текущая версия: $(node -v)"
    exit 1
fi

echo "✅ Node.js $(node -v) обнаружен"

# Создаем .env файл если его нет
if [ ! -f ".env" ]; then
    echo "📝 Создание .env файла..."
    cp .env.example .env
    echo "⚠️  Пожалуйста, отредактируйте .env файл и добавьте ваши Twitch API ключи!"
fi

# Установка зависимостей
echo "📦 Установка зависимостей..."
npm install

echo "📦 Установка зависимостей клиента..."
cd client && npm install && cd ..

# Создание папки для базы данных
echo "🗄️ Создание папки для базы данных..."
mkdir -p data

echo ""
echo "✅ Настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Отредактируйте .env файл и добавьте ваши Twitch API ключи"
echo "2. При необходимости измените список стримеров в переменной STREAMERS"
echo "3. Запустите приложение командой: npm run dev"
echo ""
echo "🔗 После запуска приложение будет доступно по адресу: http://localhost:3000"
echo "🔗 API сервер будет работать по адресу: http://localhost:3001"