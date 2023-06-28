package blacklisted_tokens

import "context"

type BlacklistedTokensRepository interface {
	Exists(ctx context.Context, token string) (bool, error)
	Add(ctx context.Context, token string, expiry int64) error
}
