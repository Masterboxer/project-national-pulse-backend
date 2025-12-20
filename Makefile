MIGRATIONS_DIR=internal/database/migrations
DB_URL=postgres://postgres:postgres@localhost:5435/journal?sslmode=disable

.PHONY: migrate-up migrate-down migrate-version migrate-force migrate-create

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(v)

MIGRATION_NAME := $(word 2, $(MAKECMDGOALS))

migrate-create:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Please specify migration name"; \
		echo "Usage: make migrate-create add_posts"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATION_NAME)

%:
	@:
