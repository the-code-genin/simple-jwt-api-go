package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
)

type Client struct {
	*redis.Client

	ttl time.Duration

	namespace string
}

// New is a client constructor.
func NewClient(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:        cfg.Host,
		Password:    cfg.Password,
		DialTimeout: 15 * time.Second,
		MaxRetries:  10,
	})

	if err := c.Ping(ctx).Err(); err != nil {
		_ = c.Close()
		return nil, err
	}

	client := &Client{
		Client:    c,
		ttl:       defaultExpirationTime,
		namespace: cfg.Prefix,
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	key = fmt.Sprintf("%s:%s", c.namespace, key)
	return c.Client.Set(ctx, key, value, duration).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, redisValue interface{}, ttl time.Duration) (bool, error) {
	key = fmt.Sprintf("%s:%s", c.namespace, key)
	return c.Client.SetNX(ctx, key, redisValue, ttl).Result()
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	key = fmt.Sprintf("%s:%s", c.namespace, key)
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Delete(ctx context.Context, key string) (int64, error) {
	key = fmt.Sprintf("%s:%s", c.namespace, key)
	return c.Client.Del(ctx, key).Result()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	key = fmt.Sprintf("%s:%s", c.namespace, key)
	i, err := c.Client.Exists(ctx, key).Result()
	return i >= 1, err
}
