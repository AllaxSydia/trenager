# Мультистадийная сборка
FROM golang:1.24-alpine AS backend

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY backend/ .
RUN go mod download
RUN go build -o main ./cmd/server

FROM node:20-alpine AS frontend

WORKDIR /app
COPY frontend/ .
RUN npm ci
RUN npm run build

# Финальный образ
FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Копируем бэкенд
COPY --from=backend /app/main .

# Копируем фронтенд
COPY --from=frontend /app/dist ./static

# Создаем пользователя
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

CMD ["./main"]