package blacklisted_tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

type BlacklistedTokensRepository struct {
	config *internal.Config
	client *redis.Client
}

// Check if the token has been blacklisted
func (tokens *BlacklistedTokensRepository) Exists(token string) (bool, error) {
	key := internal.RedisKey(tokens.config, fmt.Sprintf("blacklisted_tokens:%s", token))
	_, err := tokens.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// Blacklist a token
func (tokens *BlacklistedTokensRepository) Add(token string, expiry int64) error {
	key := internal.RedisKey(tokens.config, fmt.Sprintf("blacklisted_tokens:%s", token))
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

func NewBlacklistedTokens(config *internal.Config, client *redis.Client) BlacklistedTokens {
	return &BlacklistedTokensRepository{config, client}
}
