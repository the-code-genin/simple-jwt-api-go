.PHONY: migrateup
migrateup:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up

.PHONY: migratedown
migratedown:
	migrate -path ./database/migrations -database "postgres://postgres:password@localhost/test?sslmode=disable" up