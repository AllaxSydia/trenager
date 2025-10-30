# Root Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
# Копируем файлы зависимостей из backend
COPY backend/go.mod backend/go.sum ./

# Скачиваем зависимости
RUN go mod download

# Устанавливаем компиляторы
RUN apk add --no-cache \
    gcc \
    g++ \
    musl-dev \
    python3 \
    nodejs \
    npm \
    openjdk17 \
    openjdk17-jre

# Копируем весь backend исходный код
COPY backend/ ./

# Собираем приложение
RUN go build -o main ./cmd/server

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]