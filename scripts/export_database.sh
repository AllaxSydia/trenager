#!/bin/bash
# Скрипт для экспорта базы данных

echo "Экспорт базы данных..."
docker exec db pg_dump -U postgres trenager > database_dump_$(date +%Y%m%d_%H%M%S).sql
echo "✅ Экспорт завершен!"

