package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/the-code-genin/simple-jwt-api-go/internal/config"
)

// Connect to postgres server
func ConnectToPostgres(config *config.Config) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), config.DB.URL)
}
