package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/the-code-genin/simple-jwt-api-go/common/constants"
)

// Config stores a cache of configuration values.
type Config struct {
	Environment constants.ENV `envconfig:"ENV"`
	Port        int           `envconfig:"HTTP_PORT"`

	JWT   JWTConfig
	DB    DatabaseConfig
	Redis RedisConfig
}

func (c *Config) IsProduction() bool {
	return c.Environment == constants.ENVProd
}

type JWTConfig struct {
	Key string `envconfig:"JWT_KEY"`
	Exp int    `envconfig:"JWT_EXP"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Password string `envconfig:"REDIS_PASSWORD"`
	Prefix   string `envconfig:"REDIS_PREFIX"`
}

type DatabaseConfig struct {
	URL string `envconfig:"DATABASE_URL"`
}

// Load a new config
func LoadConfig() (*Config, error) {
	// Load env variables
	if _, err := os.Stat(".env"); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			return nil, err
		}
	}

	// Parse config data
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
