#!/bin/bash

if [ -z "$1" ]; then
    echo "Использование: ./logs.sh <сервис>"
    echo "Доступные сервисы: AuthService, TaskService, GradingService, ExecutionService, AIService, AnalyticsService, APIGateway"
    exit 1
fi

service=$1
log_file="logs/${service}.log"

if [ -f "$log_file" ]; then
    tail -f "$log_file"
else
    echo "Лог файл не найден: $log_file"
    exit 1
fi