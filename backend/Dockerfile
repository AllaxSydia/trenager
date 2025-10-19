FROM golang:1.24-alpine

WORKDIR /app

# Копируем только нужные файлы
COPY go.mod go.sum ./
RUN go mod download

# После копирования статики
COPY --from=frontend /app/dist ./static
RUN ls -la ./static/ && echo "✅ Static files copied"

COPY . .
RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]