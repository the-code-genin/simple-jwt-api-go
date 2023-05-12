.PHONY: run
run:
	go run .

.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" down

.PHONY: fmt
fmt:
	go fmt . && go fmt ./api/** && go fmt ./database/** && go fmt ./domain/** && go fmt ./internal/** && go fmt ./services/**

.PHONY: lint
lint: fmt
	golangci-lint run
