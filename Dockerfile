# Указываем какую часть собирать
ARG TARGET=backend

# Бэкенд
FROM golang:1.22-alpine AS backend-build  # ← ИЗМЕНИ НА 1.22
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Фронтенд  
FROM node:20-alpine AS frontend-build      # ← ИЗМЕНИ НА 20
WORKDIR /app
COPY frontend/ .
RUN npm ci
RUN npm run build

# Финальный образ
FROM alpine:latest

WORKDIR /app
COPY --from=backend-build /app/main ./main
COPY --from=frontend-build /app/dist ./static

EXPOSE 8080
CMD ["./main"]