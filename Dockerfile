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

# Final stage with main compilers
FROM alpine:latest

# Install main compilers and runtimes
RUN apk update && apk --no-cache add \
    ca-certificates \
    python3 \
    nodejs \
    g++ \
    # Для Go добавляем необходимые библиотеки
    libc6-compat

WORKDIR /root/

# Copy backend binary
COPY --from=backend /app/backend/main .
# Copy frontend static files
COPY --from=frontend /app/frontend/dist ./static

EXPOSE 8080
CMD ["./main"]