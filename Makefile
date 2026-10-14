.PHONY: all help test build run-backend run-frontend docker-up docker-down

help:
	@echo "🚀 URL-4-U Command Center"
	@echo "========================="
	@echo "make run-backend   - Start the Go backend API (port 8080)"
	@echo "make run-frontend  - Start the React/Vite development server (port 5173)"
	@echo "make test          - Run backend unit tests"
	@echo "make build         - Compile backend binary and build frontend bundle"
	@echo "make docker-up     - Run entire stack using Docker Compose"
	@echo "make docker-down   - Stop Docker Compose containers"

test:
	cd backend && go test -v ./...

build:
	cd backend && go build -o bin/api ./cmd/api
	cd frontend && npm run build

run-backend:
	cd backend && go run ./cmd/api/main.go

run-frontend:
	cd frontend && npm run dev

docker-up:
	docker compose up --build

docker-down:
	docker compose down
