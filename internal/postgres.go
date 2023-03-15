package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Connect to postgres server
func connectToPostgres(config *Config) (*pgx.Conn, error) {
	dbUrl, err := config.GetDBURL()
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
