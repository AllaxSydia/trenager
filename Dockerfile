# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бэкенд
COPY --from=backend /app/main .

# Копируем фронтенд
COPY --from=frontend /app/dist ./static

# Исправляем права
RUN chmod -R 755 ./static

# Запускаем от root (проще для статики)
EXPOSE 10000

CMD ["./main"]