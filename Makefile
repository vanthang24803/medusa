.PHONY: dev migrate build test tidy up down run

# Start infra (postgres, redis, minio, mailhog)
up:
	docker compose up -d

down:
	docker compose down

# Run migrations (scans modules/*/migrations)
migrate:
	go run ./apps/api/cmd/migrate

# Run API server (requires infra up + migrate beforehand)
run:
	go run ./apps/api

# Dev: infra + migrate + run
dev: up
	@sleep 3
	$(MAKE) migrate
	$(MAKE) run

build:
	go build -o bin/api ./apps/api
	go build -o bin/migrate ./apps/api/cmd/migrate

test:
	go test ./...

tidy:
	go mod tidy
