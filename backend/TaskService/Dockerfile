FROM golang:1.25-alpine

WORKDIR /app

# Копируем и собираем
COPY . .

RUN go mod download
RUN go build -o service cmd/main.go

EXPOSE 0

CMD ["./service"]
