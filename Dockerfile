# Мультистадийная сборка
FROM golang:1.24-alpine AS backend

WORKDIR /app

# Копируем и собираем бэкенд
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM node:20-alpine AS frontend

WORKDIR /app

# Копируем и собираем фронтенд
COPY frontend/ .
RUN npm ci
RUN npm run build

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бэкенд
COPY --from=backend /app/main .

# Копируем фронтенд (правильный путь!)
COPY --from=frontend /app/dist ./static

# Проверяем что файлы есть
RUN ls -la ./static/

# Создаем пользователя
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

CMD ["./main"]