#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}📊 Статус сервисов:${NC}\n"

services=("AuthService" "TaskService" "GradingService" "ExecutionService" "AIService" "AnalyticsService" "APIGateway")
ports=(50051 50052 50053 50054 50055 50056 8080)

for i in "${!services[@]}"; do
    service="${services[$i]}"
    port="${ports[$i]}"
    
    if netstat -ano 2>/dev/null | grep -q ":$port.*LISTEN"; then
        echo -e "${GREEN}✅ $service работает на порту $port${NC}"
    else
        echo -e "${RED}❌ $service НЕ работает на порту $port${NC}"
    fi
done

echo ""
echo -e "${YELLOW}🐳 Docker контейнеры:${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || echo "Docker не запущен"