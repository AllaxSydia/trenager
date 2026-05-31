#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Запуск всех микросервисов...${NC}\n"

# Создаем папки для логов и PID
mkdir -p logs pids

# Проверка и запуск PostgreSQL
if ! docker ps | grep -q postgres; then
    echo -e "${YELLOW}📦 Запуск PostgreSQL...${NC}"
    docker rm -f postgres 2>/dev/null
    docker run -d \
        --name postgres \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=postgres \
        -p 5432:5432 \
        postgres:15
    sleep 5
fi

# Создание баз данных
echo -e "${YELLOW}📦 Создание баз данных...${NC}"
docker exec -i postgres psql -U postgres << 'EOF' 2>/dev/null
CREATE DATABASE authdb;
CREATE DATABASE tasksdb;
CREATE DATABASE gradingdb;
CREATE DATABASE analyticsdb;

CREATE USER authuser WITH PASSWORD 'authpass123';
CREATE USER taskuser WITH PASSWORD 'taskpass123';
CREATE USER gradinguser WITH PASSWORD 'gradingpass123';
CREATE USER analyticsuser WITH PASSWORD 'analyticspass123';

GRANT ALL PRIVILEGES ON DATABASE authdb TO authuser;
GRANT ALL PRIVILEGES ON DATABASE tasksdb TO taskuser;
GRANT ALL PRIVILEGES ON DATABASE gradingdb TO gradinguser;
GRANT ALL PRIVILEGES ON DATABASE analyticsdb TO analyticsuser;
EOF

echo -e "${GREEN}✅ Базы данных готовы${NC}\n"

# Очищаем старые PID файлы
rm -f pids/*.pid

# Функция для запуска сервиса
start_service() {
    local name=$1
    local dir=$2
    local env=$3
    
    echo -e "${YELLOW}🔧 Запуск $name...${NC}"
    
    local current_dir=$(pwd)
    cd "$dir"
    
    if [ -n "$env" ]; then
        export $env
    fi
    
    go run cmd/main.go > "$current_dir/logs/${name}.log" 2>&1 &
    local pid=$!
    echo $pid > "$current_dir/pids/${name}.pid"
    
    cd "$current_dir"
    echo -e "${GREEN}✅ $name запущен (PID: $pid)${NC}"
}

# Запускаем сервисы
start_service "AuthService" "backend/AuthService" "DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=authdb JWT_SECRET=mysecret"
sleep 2

start_service "TaskService" "backend/TaskService" "DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=tasksdb"
sleep 2

start_service "GradingService" "backend/GradingService" "DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=gradingdb"
sleep 2

start_service "ExecutionService" "backend/ExecutionService" ""
sleep 2

start_service "AIService" "backend/AIService" ""
sleep 2

start_service "AnalyticsService" "backend/AnalyticsService" "DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=analyticsdb"
sleep 2

start_service "APIGateway" "api-gateway" ""

echo -e "\n${GREEN}🎉 Все сервисы запущены!${NC}\n"
echo -e "${YELLOW}📊 Статус:${NC}"
echo "  AuthService:       http://localhost:50051 (gRPC)"
echo "  TaskService:       http://localhost:50052 (gRPC)"
echo "  GradingService:    http://localhost:50053 (gRPC)"
echo "  ExecutionService:  http://localhost:50054 (gRPC)"
echo "  AIService:         http://localhost:50055 (gRPC)"
echo "  AnalyticsService:  http://localhost:50056 (gRPC)"
echo "  API Gateway:       http://localhost:8080 (HTTP)"
echo ""
echo -e "${YELLOW}📝 Логи:${NC}"
echo "  tail -f logs/AuthService.log"
echo "  tail -f logs/TaskService.log"
echo ""
echo -e "${RED}🛑 Остановка: ./stop-all.sh${NC}"