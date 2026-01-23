.PHONY: up down restart logs ps test build psql help

# Start the application
up:
	docker compose up -d

# Stop the application
down:
	docker compose down

# Restart the application
restart: down up

# Show logs
logs:
	docker compose logs -f

# List containers
ps:
	docker compose ps

# Build containers
build:
	docker compose build

# Run Go tests in the backend container
test:
	docker compose exec server go test ./...

# Access PostgreSQL database
psql:
	docker compose exec postgres psql -U gouser -d godb

# Help command
help:
	@echo "Available commands:"
	@echo "  make up      - Start everything in detached mode"
	@echo "  make down    - Stop and remove all containers"
	@echo "  make restart - Restart all containers"
	@echo "  make logs    - Tail logs for all containers"
	@echo "  make ps      - List running containers"
	@echo "  make build   - Rebuild application images"
	@echo "  make test    - Run tests inside the server container"
	@echo "  make psql    - Enter the database shell"
