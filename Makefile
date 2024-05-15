run: build
	@bin/astronaut-api

build:
	@go build -o bin/astronaut-api cmd/api/main.go

test:
	@go test ./... -v

