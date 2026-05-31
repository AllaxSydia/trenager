#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${RED}🛑 Остановка всех сервисов...${NC}"

# Убиваем процессы по портам
for port in 50051 50052 50053 50054 50055 50056 8080; do
    pid=$(netstat -ano 2>/dev/null | grep ":$port" | grep LISTEN | awk '{print $5}' | cut -d: -f2 | head -1)
    if [ -n "$pid" ]; then
        taskkill //F //PID $pid 2>/dev/null || kill -9 $pid 2>/dev/null
        echo -e "${GREEN}✅ Убит процесс на порту $port (PID: $pid)${NC}"
    fi
done

# Останавливаем через PID файлы
if [ -d "pids" ]; then
    for pid_file in pids/*.pid; do
        if [ -f "$pid_file" ]; then
            pid=$(cat "$pid_file")
            if kill -0 $pid 2>/dev/null; then
                kill -9 $pid 2>/dev/null
                echo -e "${GREEN}✅ Остановлен процесс (PID: $pid)${NC}"
            fi
            rm -f "$pid_file"
        fi
    done
fi

# Остановка PostgreSQL
read -p "Остановить PostgreSQL? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker stop postgres 2>/dev/null && docker rm postgres 2>/dev/null
    echo -e "${GREEN}✅ PostgreSQL остановлен и удален${NC}"
fi

echo -e "${GREEN}✅ Все сервисы остановлены${NC}"