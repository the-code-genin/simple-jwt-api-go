package blacklisted_tokens

type BlacklistedTokensRepository interface {
	Exists(token string) (bool, error)
	Add(token string, expiry int64) error
}
