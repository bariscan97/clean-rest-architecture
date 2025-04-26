.PHONY: up down migrate-up migrate-down

up:
	docker-compose up -d

down:
	docker-compose down

migrate-up:
	docker run --rm -v $(CURDIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://postgres:postgres@localhost:5432/my_db?sslmode=disable up

migrate-down:
	docker run --rm -v $(CURDIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://postgres:postgres@localhost:5432/my_db?sslmode=disable down 