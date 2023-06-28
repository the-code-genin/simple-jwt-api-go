.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" down

.PHONY: generatedocs
generatedocs:
	swag init --output services/http/docs --dir services/http,services/http/handlers,application/users -g server.go

.PHONY: generate
generate: generatedocs
	go generate -x ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run: lint
	go run ./cmd/app
