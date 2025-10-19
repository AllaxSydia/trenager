# Backend build stage
FROM golang:1.24-alpine AS backend
WORKDIR /app/backend
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Frontend build stage  
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Final stage with ALL compilers
FROM alpine:latest

# Install all compilers and runtimes
RUN apk --no-cache add \
    ca-certificates \
    python3 \
    nodejs \
    go \
    g++

WORKDIR /root/

# Copy backend binary
COPY --from=backend /app/backend/main .
# Copy frontend static files
COPY --from=frontend /app/frontend/dist ./static

EXPOSE 8080
CMD ["./main"]