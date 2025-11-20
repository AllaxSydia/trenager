FROM golang:1.24-alpine

WORKDIR /app

# Копируем workspace файлы
COPY go.work ./
COPY go.mod ./

# Копируем backend модуль
COPY backend/ ./backend/

# Скачиваем зависимости через workspace
RUN go work sync

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

# Собираем приложение через workspace
RUN go build -o main ./backend

ENV PORT=8080
EXPOSE 8080

CMD ["./main"]