.PHONY: help build run test clean docker-build docker-up docker-down fmt lint

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make deps         - Download dependencies"

deps:
	go mod download
	go mod tidy

build:
	CGO_ENABLED=0 go build -o app cmd/main.go

run: build
	./app

test:
	go test -v ./...

clean:
	rm -f app
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

docker-build:
	docker build -t expense-tracker:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

db-shell:
	docker-compose exec postgres psql -U postgres -d expense_tracker
