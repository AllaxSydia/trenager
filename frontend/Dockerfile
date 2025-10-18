FROM node:20-alpine AS builder

WORKDIR /app

COPY frontend/ .
RUN npm ci
RUN npm run build

FROM nginx:alpine

# Копируем собранное приложение
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфиг nginx
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]