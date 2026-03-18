.PHONY: dev build run test test-integration test-e2e test-all test-cover lint vet fmt fmt-check \
	migrate-up migrate-down migrate-create sqlc seed sync docker-up docker-down docker-rebuild

# Development
dev:
	air

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

run: build
	./bin/server

# Testing
test:
	go test ./... -short -count=1

test-integration:
	go test ./... -run Integration -count=1

test-e2e:
	go test ./tests/e2e/ -v -count=1 -timeout=120s

test-all:
	go test ./... -count=1

test-cover:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Code quality
lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted:" && gofmt -l . && exit 1)

# Database
migrate-up:
	migrate -path migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "$${DATABASE_URL}" down 1

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=add_tags" && exit 1)
	migrate create -ext sql -dir migrations -seq $(name)

sqlc:
	sqlc generate

seed:
	@for f in seeds/*.sql; do \
		echo "Running $$f..."; \
		psql "$${DATABASE_URL}" -f "$$f"; \
	done

sync:
	curl -X POST http://localhost:$${PORT:-8080}/api/v1/admin/sync

# Docker
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose up -d --build --force-recreate app
