FROM golang:1.24-alpine

WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk add --no-cache gcc musl-dev

# Копируем исходный код
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o main ./cmd/server

# Создаем непривилегированного пользователя
RUN adduser -D -s /bin/sh appuser
USER appuser

# Экспонируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]