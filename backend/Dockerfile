FROM golang:1.24-alpine

WORKDIR /app

# Сначала копируем только mod файлы для лучшего кэширования
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

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/server

# Для Railway важно слушать на 0.0.0.0
ENV PORT=8080
EXPOSE 8080

CMD ["./main"]