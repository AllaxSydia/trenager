FROM golang:1.21-alpine AS builder

# Устанавливаем protoc и утилиты
RUN apk add --no-cache git protoc protobuf-dev

# Устанавливаем Go плагины для protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /app

# Копируем proto файлы
COPY proto ./proto

# Генерируем protobuf код
RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/auth.proto

# ... остальной Dockerfile