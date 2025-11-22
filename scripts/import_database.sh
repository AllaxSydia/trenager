#!/bin/bash
# Скрипт для импорта базы данных
# Использование: ./import_database.sh database_dump.sql

if [ -z "$1" ]; then
    echo "❌ Укажите файл для импорта"
    echo "Использование: ./import_database.sh database_dump.sql"
    exit 1
fi

if [ ! -f "$1" ]; then
    echo "❌ Файл $1 не найден"
    exit 1
fi

echo "Импорт базы данных из $1..."
docker exec -i db psql -U postgres trenager < "$1"
echo "✅ Импорт завершен!"

