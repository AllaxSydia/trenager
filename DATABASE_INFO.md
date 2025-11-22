# Информация о базе данных

## Где хранятся данные

База данных PostgreSQL хранится в **Docker volume** локально на вашем компьютере.

### Расположение данных

**Windows:**
```
\\wsl$\docker-desktop-data\data\docker\volumes\trenager_db_data\_data
```
или через Docker Desktop можно найти в настройках Volumes.

**Linux/Mac:**
```
/var/lib/docker/volumes/trenager_db_data/_data
```

### Текущая конфигурация

- **Volume name:** `trenager_db_data`
- **Тип:** Docker named volume (локальное хранилище)
- **Содержимое:** Все данные PostgreSQL (таблицы, пользователи, индексы)

## Важно понимать

### Локальное хранилище
- ✅ Данные хранятся **локально на вашем компьютере**
- ✅ Данные **сохраняются** даже после остановки контейнеров
- ✅ Данные **удаляются** только при удалении volume (`docker volume rm trenager_db_data`)

### Для команды разработчиков

**Проблема:** У каждого разработчика своя локальная база данных с разными данными.

**Решения:**

#### 1. Общая база данных (рекомендуется для продакшена)
- Использовать внешнюю базу данных (PostgreSQL на сервере)
- Все разработчики подключаются к одной базе
- Изменения видны всем сразу

#### 2. Синхронизация через SQL дампы
- Экспортировать данные: `docker exec db pg_dump -U postgres trenager > dump.sql`
- Импортировать данные: `docker exec -i db psql -U postgres trenager < dump.sql`

#### 3. Использование миграций
- Создавать SQL скрипты для инициализации базы
- Каждый разработчик запускает миграции при первом запуске

## Полезные команды

### Просмотр данных
```bash
# Подключиться к базе данных
docker exec -it db psql -U postgres -d trenager

# Посмотреть всех пользователей
docker exec -i db psql -U postgres -d trenager -c "SELECT id, username, email, role FROM users;"

# Экспортировать данные
docker exec db pg_dump -U postgres trenager > backup.sql

# Импортировать данные
docker exec -i db psql -U postgres trenager < backup.sql
```

### Управление volume
```bash
# Посмотреть все volumes
docker volume ls

# Посмотреть информацию о volume
docker volume inspect trenager_db_data

# Удалить volume (⚠️ удалит все данные!)
docker volume rm trenager_db_data

# Создать backup volume
docker run --rm -v trenager_db_data:/data -v $(pwd):/backup alpine tar czf /backup/db_backup.tar.gz /data
```

## Рекомендации для команды

1. **Для разработки:** Каждый разработчик использует свою локальную базу
2. **Для тестирования:** Использовать общую тестовую базу данных
3. **Для продакшена:** Использовать внешнюю базу данных на сервере

## Настройка общей базы данных

Если нужно, чтобы все разработчики использовали одну базу данных, измените `docker-compose.yml`:

```yaml
db:
  image: postgres:15
  container_name: db
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: trenager
  # Убрать volumes для локального хранения
  # volumes:
  #   - db_data:/var/lib/postgresql/data
  # Вместо этого использовать внешнюю БД
  # Или подключиться к существующей БД через переменные окружения
```

И в `backend` изменить переменные окружения:
```yaml
environment:
  - DB_HOST=your-database-server.com  # Внешний сервер БД
  - DB_PORT=5432
  - DB_USER=postgres
  - DB_PASSWORD=your_password
  - DB_NAME=trenager
```

