include .env
MIGRATION_PATH=./migrations

.PHONY: migrate-create migrate-up migrate-down

run:
	@go run ./cmd/server
	
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(name) 

migrate-up:
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1