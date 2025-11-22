# Root Makefile to orchestrate Docker Compose
# Acts as a simple wrapper around docker compose commands.

.PHONY: all build up down clean

all: build

# Build all services defined in docker-compose.yml
build:
	@echo "Building services..."
	docker compose build

# Run the services in detached mode
up:
	@echo "Starting services..."
	docker compose up -d

# Stop and remove the services
down:
	@echo "Stopping services..."
	docker compose down

# Clean up Docker images (use with caution)
clean:
	@echo "Cleaning up..."
	docker compose down --rmi all -v