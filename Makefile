.PHONY: run
run:
	go run ./cmd/app

.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/go_jwt_api?sslmode=disable" down

.PHONY: fmt
fmt:
	go fmt ./api/** && go fmt ./application/** && go fmt ./common/** && go fmt ./database/ && go fmt ./domain/**

.PHONY: lint
lint: fmt
	golangci-lint run
