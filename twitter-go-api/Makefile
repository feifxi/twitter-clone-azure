COMPOSE_FILE ?= ../docker-compose.yml
DB_SERVICE ?= database
DB_URL ?= postgres://root:rootpass@localhost:5432/twitter_db?sslmode=disable

.PHONY: db-up db-down db-restart db-logs db-ps db-shell db-wait migrateup migratedown sqlc test run

db-up:
	docker compose -f $(COMPOSE_FILE) up -d $(DB_SERVICE)

db-down:
	docker compose -f $(COMPOSE_FILE) down

db-restart: db-down db-up db-wait

db-logs:
	docker compose -f $(COMPOSE_FILE) logs -f $(DB_SERVICE)

db-ps:
	docker compose -f $(COMPOSE_FILE) ps

db-shell:
	docker compose -f $(COMPOSE_FILE) exec $(DB_SERVICE) psql -U root -d twitter_db

db-wait:
	@until docker compose -f $(COMPOSE_FILE) exec -T $(DB_SERVICE) pg_isready -U root -d twitter_db >/dev/null 2>&1; do \
		echo "waiting for postgres..."; \
		sleep 1; \
	done

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	GOCACHE=/tmp/go-build-cache go test -v -cover ./...

run:
	go run cmd/api/main.go
