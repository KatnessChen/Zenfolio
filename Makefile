.PHONY: dev-up dev-down dev-build dev-logs dev-test

# 構建並啟動所有服務
dev-up:
	docker-compose up -d --build

# 只構建不啟動
dev-build:
	docker-compose build

# 停止並移除所有服務
dev-down:
	docker-compose down

# 查看服務日誌
dev-logs:
	docker-compose logs -f

# 查看特定服務日誌
dev-logs-backend:
	docker-compose logs -f backend

dev-logs-price:
	docker-compose logs -f price-service

# 運行測試
dev-test-backend:
	cd backend && go test -v ./...

dev-test-price:
	cd price_service && go test -v ./...

# 進入容器 shell
dev-shell-backend:
	docker-compose exec backend sh

dev-shell-price:
	docker-compose exec price-service sh

# 重啟特定服務
dev-restart-backend:
	docker-compose restart backend

dev-restart-price:
	docker-compose restart price-service

# 清理所有容器和卷
dev-clean:
	docker-compose down -v

# 顯示容器狀態
dev-ps:
	docker-compose ps

# 顯示容器資源使用情況
dev-stats:
	docker stats
