FROM golang:1.24-alpine

WORKDIR /app

# Копируем mod файлы
COPY go.mod go.sum ./

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

# Копируем исходный код из backend
COPY backend/ ./backend/

# Собираем приложение
RUN go build -o main ./backend/cmd/server

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]