package redis

import (
	"fmt"

	"github.com/the-code-genin/simple-jwt-api-go/common/config"
)

// Prefixes the key with the app redis key for namespacing
func RedisKey(config *config.Config, key string) string {
	return fmt.Sprintf("%s:%s", config.Redis.Prefix, key)
}
