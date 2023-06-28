package blacklisted_tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/the-code-genin/simple-jwt-api-go/common/redis"
)

type blacklistedTokensRepository struct {
	client *redis.Client
}

func (tokens *blacklistedTokensRepository) Exists(ctx context.Context, token string) (bool, error) {
	res, err := tokens.client.Exists(ctx, fmt.Sprintf("blacklisted_tokens:%s", token))
	if err != nil {
		return false, err
	}
	return res, nil
}

func (tokens *blacklistedTokensRepository) Add(ctx context.Context, token string, expiry int64) error {
	err := tokens.client.Set(
		ctx,
		fmt.Sprintf("blacklisted_tokens:%s", token),
		expiry,
		time.Until(time.Unix(expiry, 0)),
	)
	return err
}

func NewBlacklistedTokensRepository(client *redis.Client) BlacklistedTokensRepository {
	return &blacklistedTokensRepository{client}
}
