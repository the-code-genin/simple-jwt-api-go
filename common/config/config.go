package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config stores a cache of configuration values.
type Config struct {
	Env   string `envconfig:"env"`
	Port  int    `envconfig:"http_port"`
	JWT   JWTConfig
	DB    DatabaseConfig
	Redis RedisConfig
}

type JWTConfig struct {
	Key string `envconfig:"jwt_key"`
	Exp int    `envconfig:"jwt_exp"`
}

type RedisConfig struct {
	Host     string `envconfig:"redis_host"`
	Password string `envconfig:"redis_password"`
	Prefix   string `envconfig:"redis_prefix"`
}

type DatabaseConfig struct {
	URL string `envconfig:"database_url"`
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
