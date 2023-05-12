package repositories

type BlacklistedTokensRepository interface {
	Exists(token string) (bool, error)
	Add(token string, expiry int64) error
}
