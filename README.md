# Get Rates USDT GRPC Service

GRPC-сервис для получения курса USDT с биржи Garantex с сохранением данных в PostgreSQL.

## Основные функции
- Получение курса USDT через GRPC метод GetRates
- Сохранение данных с меткой времени в PostgreSQL
- Healthcheck для проверки работоспособности
- Мониторинг через Prometheus
- Трассировка запросов через OpenTelemetry
- Graceful shutdown
- Автоматические миграции БД

## Требования
- Go 1.22+
- Docker и Docker Compose
- Protoc (для генерации кода из .proto)
- PostgreSQL (запускается через Docker)
- Prometheus (запускается через Docker)
- Jaeger (запускается через Docker)


### Установка
```bash
git clone https://github.com/dgoodandrun/get-rates-usdt-grpc-service.git
cd get-rates-usdt-grpc-service
go mod download
make protoc
```
### Настройка
Параметры для настройки сервиса в корне проекта файл .env
```bash
# App name
APPNAME=getRates
# GRPC
PORT=50051
# Prometheus
METRICS_PORT=9090
# PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=rates
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
# Jaeger
JAEGER_WEB_UI_PORT=16686
JAEGER_COLLECTOR_PORT=14268
# Garantex API
GARANTEX_API_URL=https://garantex.org/api/v2/depth?market=%s
GARANTEX_API_URL_MARKET=btcusdt
```
### Запуск
```bash
docker-compose up --build
```
### Запуск тестов
```bash
make test
```
### Запуск линтера
```bash
make lint
```