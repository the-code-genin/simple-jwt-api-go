.PHONY: run
run:
	go run .

.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" down

.PHONY: fmt
fmt:
	go fmt main.go && go fmt ./api && go fmt ./database/blacklisted_tokens && go fmt ./database/users && go fmt ./domain && go fmt ./internal && go fmt ./services

.PHONY: lint
lint: fmt
	golangci-lint run
