# Мультистадийная сборка
FROM golang:1.22-alpine AS backend

WORKDIR /app

# Копируем и собираем бэкенд
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM node:20-alpine AS frontend

WORKDIR /app

# Копируем и собираем фронтенд
COPY frontend/ ./frontend/
WORKDIR /app/frontend
RUN npm ci
RUN npm run build

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бэкенд
COPY --from=backend /app/backend/main .

# Копируем фронтенд
COPY --from=frontend /app/frontend/dist ./static

# Создаем пользователя
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

CMD ["./main"]