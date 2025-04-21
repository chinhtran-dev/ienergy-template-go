.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create swagger-init swagger-build build run test test-coverage install-tools test-all test-unit test-integration test-http lint

install-tools:
	@echo "Installing required tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo "Tools installed successfully!"

migrate-up:
	go run cmd/migrate/main.go -action up

migrate-down:
	go run cmd/migrate/main.go -action down

migrate-force:
	go run cmd/migrate/main.go -action force $(version)

migrate-version:
	go run cmd/migrate/main.go -action version

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

swagger-init:
	swag init -g cmd/app/main.go -o docs/swagger

swagger-build:
	swag fmt
	swag init -g cmd/app/main.go -o docs/swagger

build:
	go build -o bin/app ./cmd/app/main.go

run:
	bin/app api

test:
	go test -v ./...

test-all: test-unit test-integration test-http

test-unit:
	go test -v ./internal/... -tags=unit

test-integration:
	go test -v ./tests/integration/... -tags=integration

test-http:
	go test -v ./internal/http/... -tags=http

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...