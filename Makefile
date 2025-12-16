MIGRATIONS_DIR=internal/database/migrations
DB_URL=postgres://postgres:postgres@localhost:5435/journal?sslmode=disable

.PHONY: migrate-up migrate-down migrate-version migrate-create

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(v)

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
