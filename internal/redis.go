package internal

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect to redis server
func ConnectToRedis(config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,
	})
	return client, nil
}

// Prefixes the key with the app redis key for namespacing
func RedisKey(config Config, key string) string {
	return fmt.Sprintf("%s.%s", config.Redis.Prefix, key)
}
