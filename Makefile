.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create

migrate-up:
	go run cmd/migrate/main.go -action up

migrate-down:
	go run cmd/migrate/main.go -action down

migrate-force:
	go run cmd/migrate/main.go -action force $(version)

migrate-version:
	go run cmd/migrate/main.go -action version

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)