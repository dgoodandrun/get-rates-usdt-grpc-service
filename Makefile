# Makefile

APP_NAME = get-rates-service
PROTOGEN_DIR = protogen/golang
PROTO_FILE = get-rates/get-rates.proto
PROTO_DIR = proto
DOCKER_TAG = latest




.PHONY: protoc build test docker-build run lint

protoc:
	@echo "Generating protobuf code..."
	cd $(PROTO_DIR) && protoc \
    --go_out=../$(PROTOGEN_DIR) --go_opt=paths=source_relative \
    --go-grpc_out=../$(PROTOGEN_DIR) --go-grpc_opt=paths=source_relative \
    $(PROTO_FILE)

build: protoc
	@echo "Building application..."
	go build -o bin/$(APP_NAME) ./cmd/main.go

test:
	@echo "Running tests..."
	go test -v ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(DOCKER_TAG) .

run: build
	@echo "Starting application..."
	./bin/$(APP_NAME)

lint:
	@echo "Running linter..."
	golangci-lint run ./...