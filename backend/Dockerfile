FROM golang:1.24-alpine

WORKDIR /app

# Установка компиляторов
RUN apk add --no-cache \
    python3 \
    nodejs \
    npm \
    openjdk11 \
    g++ \
    make \
    git

# Копируем файлы зависимостей ИЗ ПАПКИ backend
COPY backend/go.mod backend/go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код бэкенда
COPY backend/ ./

# Собираем приложение
RUN go build -o main ./cmd/server

# Создаем пользователя и настраиваем права
RUN adduser -D -s /bin/sh appuser
RUN mkdir -p /tmp/code_executor && \
    chown appuser:appuser /tmp/code_executor

USER appuser

EXPOSE 8080

CMD ["./main"]