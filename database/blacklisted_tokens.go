package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/the-code-genin/simple-jwt-api-go/internal"
)

var mutex sync.Mutex

type BlacklistedTokens struct {
	ctx *internal.AppContext
}

// Check if the token has been blacklisted
func (tokens *BlacklistedTokens) Exists(token string) (bool, error) {
	client, err := tokens.ctx.GetRedisClient()
	if err != nil {
		return false, err
	}
	key, err := internal.RedisKey(tokens.ctx, fmt.Sprintf("blacklisted_tokens:%s", token))
	if err != nil {
		return false, err
	}

	_, err = client.Get(context.Background(), key).Result()
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
func (tokens *BlacklistedTokens) Add(token string, expiry int64) error {
	client, err := tokens.ctx.GetRedisClient()
	if err != nil {
		return err
	}
	key, err := internal.RedisKey(tokens.ctx, fmt.Sprintf("blacklisted_tokens:%s", token))
	if err != nil {
		return err
	}

	_, err = client.Set(context.Background(), key, expiry, time.Unix(expiry, 0).Sub(time.Now())).Result()
	if err != nil {
		return err
	}
	return err
}

func NewBlacklistedTokens(ctx *internal.AppContext) *BlacklistedTokens {
	return &BlacklistedTokens{ctx}
}
