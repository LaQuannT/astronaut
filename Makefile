include .env

run: build
	@bin/astronaut-api

build:
	@go build -o bin/astronaut-api cmd/api/main.go

test:
	@go test ./... -v

migrate-create:
	@migrate create -ext sql -dir migration/  -seq $(NAME)

migrate-up:
	@migrate -path migration/ -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" -verbose up

migrate-down:
	@migrate -path migration/ -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"  -verbose down

migrate-fix:
	@migrate -path migration/ -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" force $(VERSION)