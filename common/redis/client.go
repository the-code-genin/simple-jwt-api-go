package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
)

func NewClient(config *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,
	})
	return client, nil
}
