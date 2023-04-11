.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: sqlgenerate
sqlgenerate:
	sqlboiler psql -c ./sqlboiler.toml

.PHONY: fmt
fmt:
	go fmt main.go && go fmt ./api && go fmt ./database/repositories && go fmt ./internal