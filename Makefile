.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" down

.PHONY: fmt
fmt:
	go fmt ./api/** && go fmt ./application/** && go fmt ./common/** && go fmt ./database/** && swag fmt --dir api/http,application/users -g server.go

.PHONY: generatedocs
generatedocs: fmt
	swag init --dir api/http,application/users -g server.go

.PHONY: lint
lint: fmt generatedocs
	golangci-lint run

.PHONY: run
run: lint
	go run ./cmd/app
