.PHONY: dev-up dev-down dev-build dev-logs dev-test

# Build and start all services
dev-up:
	docker-compose up -d --build

# Build only, do not start
dev-build:
	docker-compose build

# Stop and remove all services
dev-down:
	docker-compose down

# View service logs
dev-logs:
	docker-compose logs -f

# View logs for specific service
dev-logs-backend:
	docker-compose logs -f backend

dev-logs-price:
	docker-compose logs -f price-service

# Run tests
dev-test-backend:
	cd backend && go test -v ./...

dev-test-price:
	cd price_service && go test -v ./...

# Enter container shell
dev-shell-backend:
	docker-compose exec backend sh

dev-shell-price:
	docker-compose exec price-service sh

# Restart specific service
dev-restart-backend:
	docker-compose restart backend

dev-restart-price:
	docker-compose restart price-service

# Clean all containers and volumes
dev-clean:
	docker-compose down -v

# Show container status
dev-ps:
	docker-compose ps

# Show container resource usage
dev-stats:
	docker stats
