package internal

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
