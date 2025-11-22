@echo off
REM Скрипт для импорта базы данных (Windows)
REM Использование: import_database.bat database_dump.sql

if "%1"=="" (
    echo ❌ Укажите файл для импорта
    echo Использование: import_database.bat database_dump.sql
    pause
    exit /b 1
)

if not exist "%1" (
    echo ❌ Файл %1 не найден
    pause
    exit /b 1
)

echo Импорт базы данных из %1...
docker exec -i db psql -U postgres trenager < "%1"
echo ✅ Импорт завершен!
pause

