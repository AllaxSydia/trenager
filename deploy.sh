#!/bin/bash

echo "🚀 Starting deployment..."

# Останавливаем существующие контейнеры
docker-compose down

# Собираем и запускаем
docker-compose up -d --build

echo "✅ Deployment completed!"
echo "📊 Checking services..."

# Проверяем статус
sleep 10
docker-compose ps

echo "🌐 Backend URL: http://your-server-ip:8080"
echo "🎯 Frontend URL: http://your-server-ip"
echo "📋 Health check: http://your-server-ip:8080/health"