include .env

run: build
	@bin/astronaut-data-api

build:
	@go build -o bin/astronaut-data-api cmd/api/main.go

test:
	@go test ./...

docker-compose:
	@docker-compose up -d

migration_create:
	@migrate create -ext sql -dir internal/database/migration/ -seq $(NAME)

migration_up: 
	@migrate -path internal/database/migration/ -database "postgresql://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}" -verbose up

migration_down: 
	@migrate -path internal/database/migration/ -database "postgresql://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}"  -verbose down

migration_fix: 
	@migrate -path internal/database/migration/ -database "postgresql://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}" force $(VERSION)
