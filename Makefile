include .env
MIGRATION_PATH=./migrations

.PHONY: run seed migrate-create migrate-up migrate-down

# ---- App ----

run: ## 啟動 API server
	@go run ./cmd/server

seed: ## 塞入開發用假資料
	@go run ./cmd/seed


# ---- Migration ----

migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(name) 

migrate-up:
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1
