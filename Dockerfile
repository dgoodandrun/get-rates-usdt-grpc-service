# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарный файл сервиса
RUN go build -o /app/main ./cmd

# Используем минимальный образ для запуска сервиса
FROM alpine:latest

# Устанавливаем CA сертификаты для TLS
RUN apk add --no-cache ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего шага
COPY --from=builder /app/main /app/main

# Открываем порт для gRPC
EXPOSE 50051

# Запускаем приложение
CMD ["/app/main"]