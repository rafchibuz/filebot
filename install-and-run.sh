#!/bin/bash

echo "🎬 Twitch Clips Aggregator - Полная установка и запуск"
echo "=================================================="

# Проверяем наличие Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js не установлен."
    echo "📥 Установка Node.js..."
    
    # Для Ubuntu/Debian
    if command -v apt-get &> /dev/null; then
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        sudo apt-get install -y nodejs
    # Для CentOS/RHEL/Fedora
    elif command -v yum &> /dev/null; then
        curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
        sudo yum install -y nodejs
    else
        echo "❌ Автоматическая установка Node.js не поддерживается для вашей системы."
        echo "Пожалуйста, установите Node.js вручную с https://nodejs.org/"
        exit 1
    fi
fi

echo "✅ Node.js $(node -v) готов к работе"

# Создание .env файла
if [ ! -f ".env" ]; then
    echo "📝 Создание конфигурационного файла..."
    cp .env.example .env
fi

# Установка зависимостей
echo "📦 Установка зависимостей сервера..."
npm install

echo "📦 Установка зависимостей клиента..."
cd client && npm install && cd ..

# Создание папки для базы данных
mkdir -p data

echo ""
echo "✅ Установка завершена!"
echo ""
echo "⚙️  ВАЖНО: Настройте Twitch API ключи в файле .env"
echo ""
echo "📋 Инструкция по получению ключей:"
echo "1. Перейдите на https://dev.twitch.tv/console/apps"
echo "2. Нажмите 'Register Your Application'"
echo "3. Заполните форму:"
echo "   - Name: Twitch Clips Aggregator"
echo "   - OAuth Redirect URLs: http://localhost:3000"
echo "   - Category: Analytics Tool"
echo "4. Скопируйте Client ID и Client Secret в .env файл"
echo ""
echo "🎯 Настроенные стримеры: RavshanN, steel, renatko"
echo "   (можно изменить в .env файле)"
echo ""

# Проверяем настройку API ключей
if grep -q "your_client_id_here" .env; then
    echo "⚠️  Для запуска приложения необходимо настроить Twitch API ключи в .env файле"
    echo ""
    echo "Откройте файл .env и замените:"
    echo "  TWITCH_CLIENT_ID=your_client_id_here"
    echo "  TWITCH_CLIENT_SECRET=your_client_secret_here"
    echo ""
    echo "После настройки запустите: ./start.sh"
else
    echo "🚀 Запуск приложения..."
    ./start.sh
fi