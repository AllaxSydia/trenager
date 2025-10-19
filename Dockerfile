# Backend build stage
FROM golang:1.24-alpine as backend
WORKDIR /app/backend
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Frontend build stage  
FROM node:20-alpine as frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Final stage - минимальный набор
FROM alpine:latest

# Устанавливаем только Python и Node.js (C++ оставляем т.к. он работает)
RUN apk update && apk --no-cache add \
    ca-certificates \
    python3 \
    nodejs \
    g++

WORKDIR /root/

# Copy backend binary
COPY --from=backend /app/backend/main .
# Copy frontend static files
COPY --from=frontend /app/frontend/dist ./static

EXPOSE 8080
CMD ["./main"]