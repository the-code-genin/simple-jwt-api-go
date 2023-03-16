package internal

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Connect to redis server
func connectToRedis(config *Config) (*redis.Client, error) {
	redisHost, err := config.GetRedisHost()
	if err != nil {
		return nil, err
	}

	redisPassword, err := config.GetRedisPassword()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})
	return client, nil
}

// Prefixes the key with the app redis key for namespacing
func RedisKey(ctx *AppContext, key string) (string, error) {
	config := ctx.GetConfig()
	redisPrefix, err := config.GetRedisPrefix()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", redisPrefix, key), nil
}
