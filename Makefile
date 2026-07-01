SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

DB_DSN ?= mysql://root:admin@tcp(localhost:3306)/goldin
MIGRATIONS_DIR ?= migrations/mysql

.PHONY: test compose-up compose-down migrate-up migrate-down migrate-force

test:
	go test ./...

compose-up:
	docker compose up -d

compose-down:
	docker compose down

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

migrate-force:
	@if ("$(VERSION)" -eq "") { throw "VERSION is required. Usage: make migrate-force VERSION=1" }
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" force $(VERSION)
