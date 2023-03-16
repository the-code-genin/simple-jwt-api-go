package internal

import (
	"crypto/sha256"
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores a cache of configuration values.
type Config struct {
	cache map[string]string
}

// Error if configuration value is not found
var ErrorConfigNotFound = errors.New("config not found")

// Config attempts to fetch results from it's internal cache,
// if the data is not in the internal cache, it attempts to get it from the env variables,
// if the data is not set in the env variables, it returns a not found error.
func (c *Config) get(name string) (string, error) {
	// Return the value from cache if it is available
	if value, ok := c.cache[name]; ok {
		return value, nil
	}

	// Return the value from the env variables if it is available
	// The value is cached for future queries
	if value, ok := os.LookupEnv(name); ok {
		c.cache[name] = value
		return value, nil
	}

	// Value not available
	return "", ErrorConfigNotFound
}

// Get the env
func (c *Config) GetEnv() (string, error) {
	return c.get("ENV")
}

// Get the JWT SHA256-HMAC key
func (c *Config) GetJWTKey() ([]byte, error) {
	val, err := c.get("JWT_KEY")
	if err != nil {
		return nil, err
	}
	key := sha256.Sum256([]byte(val))
	return key[:], nil
}

// Get the JWT expiry period in seconds
func (c *Config) GetJWTExpiry() (int, error) {
	val, err := c.get("JWT_EXP")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

// Get the HTTP Port
func (c *Config) GetHTTPPort() (int, error) {
	val, err := c.get("HTTP_PORT")
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return port, nil
}

// Get the postgres database URL
func (c *Config) GetDBURL() (string, error) {
	return c.get("DATABASE_URL")
}

// Get the redis host
func (c *Config) GetRedisHost() (string, error) {
	return c.get("REDIS_HOST")
}

// Get the redis password
func (c *Config) GetRedisPassword() (string, error) {
	return c.get("REDIS_PASSWORD")
}

// Get the redis prefix
func (c *Config) GetRedisPrefix() (string, error) {
	return c.get("REDIS_PREFIX")
}

// Attempts to load .env variables into the config if .env file exists
func loadConfig() (config *Config) {
	config = &Config{make(map[string]string)}

	// Skip operation if .env file does not exist
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}

	// Load env variables into memory
	godotenv.Load(".env")
	return
}
