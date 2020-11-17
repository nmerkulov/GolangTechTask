include .env
export
.PHONY: run
run: migrate
	go run cmd/server/main.go
migrate:
	migrate -path cmd/server/internal/migrations -database 'postgres://postgres:secretpassword@localhost:5432/buffs?sslmode=disable' up