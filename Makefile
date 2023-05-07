.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: fmt
fmt:
	go fmt main.go && go fmt ./api && go fmt ./database/repositories && go fmt ./internal

.PHONY: lint
lint: fmt
	golangci-lint run
