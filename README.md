# Get Rates USDT GRPC Service

GRPC-сервис для получения курса USDT с биржи Garantex с сохранением данных в ClickHouse.

## Основные функции
- Получение курса USDT через GRPC метод GetRates
- Сохранение данных с меткой времени в ClickHouse
- Healthcheck для проверки работоспособности
- Мониторинг через Prometheus
- Трассировка запросов через OpenTelemetry
- Graceful shutdown
- Автоматические миграции БД

## Требования
- Go 1.22+
- Docker и Docker Compose
- Protoc (для генерации кода из .proto)
- ClickHouse (запускается через Docker)


### Установка
```bash
git clone https://github.com/dgoodandrun/get-rates-usdt-grpc-service.git
cd get-rates-usdt-grpc-service
go mod download
make protoc
```
### Настройка
Параметры для настройки сервиса в корне проекта файл .env

### Запуск
```bash
docker-compose up --build
```
### Запуск тестов
```bash
make test
```