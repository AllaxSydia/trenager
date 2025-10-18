# Мультистадийная сборка
FROM golang:1.24-alpine AS backend

WORKDIR /app

# Копируем и собираем бэкенд (из папки backend на одном уровне с Dockerfile)
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM node:20-alpine AS frontend

WORKDIR /app

# Копируем и собираем фронтенд (из папки frontend на одном уровне с Dockerfile)
COPY frontend/ ./frontend/
WORKDIR /app/frontend
RUN npm ci
RUN npm run build

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем бэкенд
COPY --from=backend /app/backend/main .

# Копируем фронтенд
COPY --from=frontend /app/frontend/dist ./static

# Исправляем права
RUN chmod -R 755 ./static

EXPOSE 10000

CMD ["./main"]