# Makefile

.PHONY: run test lint build docker-build docker-run

run:
	go run ./cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server/main.go

docker-build:
	docker build -t go-microservice-boilerplate:latest .

docker-run:
	docker compose up --build
