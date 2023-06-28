package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
)

func NewConnection(config *config.DatabaseConfig) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), config.URL)
}
