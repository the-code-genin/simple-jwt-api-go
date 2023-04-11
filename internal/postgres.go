package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Connect to postgres server
func ConnectToPostgres(config *Config) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), config.DB.URL)
}
