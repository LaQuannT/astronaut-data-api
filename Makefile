run: build
	@bin/astronaut-data-api

build:
	@go build -o bin/astronaut-data-api cmd/api/main.go

test:
	@go test ./...

docker-compose:
	@docker-compose up -d
