package database

import (
	"context"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
	"github.com/the-code-genin/simple-jwt-api-go/domain/repositories"
)

type BlacklistedTokensRepository struct {
	config *config.Config
	client *r.Client
}

func (tokens *BlacklistedTokensRepository) Exists(token string) (bool, error) {
	key := redis.RedisKey(tokens.config, fmt.Sprintf("blacklisted_tokens:%s", token))
	_, err := tokens.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == r.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (tokens *BlacklistedTokensRepository) Add(token string, expiry int64) error {
	key := redis.RedisKey(tokens.config, fmt.Sprintf("blacklisted_tokens:%s", token))
	_, err := tokens.client.Set(
		context.Background(),
		key,
		expiry,
		time.Until(time.Unix(expiry, 0)),
	).Result()
	if err != nil {
		return err
	}
	return err
}

func NewBlacklistedTokensRepository(config *config.Config, client *r.Client) repositories.BlacklistedTokensRepository {
	return &BlacklistedTokensRepository{config, client}
}
