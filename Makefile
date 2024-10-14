DATABASE_URL ?= postgres://user:password@host:port/db-name?sslmode=disable

.PHONY: migrate build run run_local down

build:
	docker compose build

run: build
	docker compose up -d

migrate: run
	docker compose exec academ_stats migrate -path db/migrations -database "$(DATABASE_URL)" up

down:
	docker compose down

run_local:
	@echo "Define local run configuration if needed"
