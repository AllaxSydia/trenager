# Корневой Dockerfile
ARG APP=backend

FROM node:18-alpine as frontend
WORKDIR /app
COPY frontend/ .
RUN npm ci
RUN npm run build

FROM golang:1.21-alpine as backend
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest as production

# Выбираем что запускать через ARG
ARG APP
WORKDIR /root/

COPY --from=backend /app/main ./main
COPY --from=frontend /app/dist ./frontend-dist

# Запускаем выбранное приложение
CMD if [ "$APP" = "backend" ]; then ./main; else cd frontend-dist && nginx -g 'daemon off;'; fi

EXPOSE 8080 3000