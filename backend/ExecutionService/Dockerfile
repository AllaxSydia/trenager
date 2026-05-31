FROM golang:1.25-alpine

# Устанавливаем runtime для разных языков
RUN apk add --no-cache python3 py3-pip nodejs openjdk21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o execution-service cmd/main.go

EXPOSE 50054

CMD ["./execution-service"]