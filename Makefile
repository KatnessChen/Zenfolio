.PHONY: dev-up dev-down dev-build dev-logs dev-test dev-lint dev-format

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
dev-logs-frontend:
	docker-compose logs -f frontend

dev-logs-backend:
	docker-compose logs -f backend

dev-logs-price:
	docker-compose logs -f price-service

# Run tests
dev-test:
	make dev-test-frontend
	make dev-test-backend
	make dev-test-price

dev-test-frontend:
	docker-compose exec frontend pnpm test

dev-test-backend:
	cd backend && go test -v ./...

dev-test-price:
	cd price_service && go test -v ./...

# Code quality
dev-lint:
	make dev-lint-frontend
	make dev-lint-backend

dev-lint-frontend:
	docker-compose exec frontend pnpm lint

dev-lint-frontend-fix:
	docker-compose exec frontend pnpm lint:fix

dev-lint-backend:
	cd backend && golangci-lint run

dev-format:
	make dev-format-frontend
	make dev-format-backend

dev-format-frontend:
	docker-compose exec frontend pnpm format

dev-format-frontend-check:
	docker-compose exec frontend pnpm format:check

dev-format-backend:
	cd backend && gofmt -w .

# Type checking
dev-type-check-frontend:
	docker-compose exec frontend pnpm type-check

# Enter container shell
dev-shell-frontend:
	docker-compose exec frontend sh

dev-shell-backend:
	docker-compose exec backend sh

dev-shell-price:
	docker-compose exec price-service sh

# Restart specific service
dev-restart-frontend:
	docker-compose restart frontend

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
