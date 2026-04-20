include .env
export 

postgres-up:
	docker compose -f docker-compose.yml up -d kinotower-postgres

postgres-down:
	docker compose -f docker-compose.yml down kinotower-postgres

migrate-create:
	docker compose -f docker-compose.yml run --rm kinotower-postgres-migrate create -ext sql -dir /migrations $(name)

migrate-up:
	docker compose -f docker-compose.yml run --rm kinotower-postgres-migrate -path=/migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@kinotower-postgres:5432/$(POSTGRES_DB)?sslmode=disable" up

migrate-down:
	docker compose -f docker-compose.yml run --rm kinotower-postgres-migrate -path=/migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@kinotower-postgres:5432/$(POSTGRES_DB)?sslmode=disable" down

kinotower-run:
	go mod tidy && go run cmd/kinotower/main.go