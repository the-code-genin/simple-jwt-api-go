package internal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type ContextKey uint8

const (
	configContextKey   ContextKey = 1
	redisContextKey    ContextKey = 2
	postgresContextKey ContextKey = 3
)

type AppContext struct {
	ctx context.Context
}

// Get configuration from the context
func (s *AppContext) GetConfig() *Config {
	config, ok := s.ctx.Value(configContextKey).(*Config)
	if !ok {
		config := loadConfig()
		s.ctx = context.WithValue(s.ctx, configContextKey, config)
		return config
	}
	return config
}

// Get redis client from the context
func (s *AppContext) GetRedisClient() (*redis.Client, error) {
	client, ok := s.ctx.Value(redisContextKey).(*redis.Client)
	if !ok {
		client, err := connectToRedis(s.GetConfig())
		if err != nil {
			return nil, err
		}
		s.ctx = context.WithValue(s.ctx, redisContextKey, client)
		return client, nil
	}
	return client, nil
}

// Get postgres connection from the context
func (s *AppContext) GetPostgresConn() (*pgx.Conn, error) {
	conn, ok := s.ctx.Value(postgresContextKey).(*pgx.Conn)
	if !ok {
		conn, err := connectToPostgres(s.GetConfig())
		if err != nil {
			return nil, err
		}
		s.ctx = context.WithValue(s.ctx, postgresContextKey, conn)
		return conn, nil
	}
	return conn, nil
}

func NewAppContext(ctx context.Context) *AppContext {
	return &AppContext{ctx}
}
