@echo off
REM Скрипт для экспорта базы данных (Windows)

echo Экспорт базы данных...
docker exec db pg_dump -U postgres trenager > database_dump_%date:~-4,4%%date:~-7,2%%date:~-10,2%_%time:~0,2%%time:~3,2%%time:~6,2%.sql
echo ✅ Экспорт завершен!
pause

