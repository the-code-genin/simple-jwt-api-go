package internal

import (
	"database/sql"
)

// Connect to postgres server
func ConnectToPostgres(config Config) (*sql.DB, error) {
	return sql.Open("postgres", config.DB.URL)
}
