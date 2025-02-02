# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарный файл сервиса
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

# Устанавливаем CA сертификаты для TLS
RUN apk add --no-cache ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из предыдущего шага
COPY --from=builder /app/main /app/main

# Открываем порты для gRPC и метрик Prometheus
EXPOSE 50051 9090

# Запускаем приложение
CMD ["/app/main"]