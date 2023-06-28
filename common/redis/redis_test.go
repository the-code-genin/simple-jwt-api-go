package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/the-code-genin/simple-jwt-api-go/common/config"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()

	conf, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, conf)

	redisClient, err := NewClient(ctx, conf.Redis)
	assert.NoError(t, err)
	assert.NotNil(t, redisClient)

	t.Run("TestRedisSetterGetter", func(t *testing.T) {
		err := redisClient.Set(ctx, "atestkey", 5838, 5*time.Second)
		assert.NoError(t, err)

		time.Sleep(3 * time.Second)

		value, err := redisClient.Get(ctx, "atestkey")
		assert.NoError(t, err)
		assert.Equal(t, "5838", value)

		time.Sleep(3 * time.Second)

		value, err = redisClient.Get(ctx, "atestkey")
		assert.Error(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("TestRedisDeleteExists", func(t *testing.T) {
		err := redisClient.Set(ctx, "btestkey", 5838, 5*time.Second)
		assert.NoError(t, err)

		time.Sleep(3 * time.Second)

		value, err := redisClient.Exists(ctx, "btestkey")
		assert.NoError(t, err)
		assert.True(t, value)

		count, err := redisClient.Delete(ctx, "btestkey")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		value, err = redisClient.Exists(ctx, "btestkey")
		assert.NoError(t, err)
		assert.False(t, value)

		count, err = redisClient.Delete(ctx, "btestkey")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("TestRedisSetNX", func(t *testing.T) {
		_, err := redisClient.Delete(ctx, "ctestkey")
		assert.NoError(t, err)

		res, err := redisClient.SetNX(ctx, "ctestkey", 5838, 5*time.Second)
		assert.NoError(t, err)
		assert.True(t, res)

		time.Sleep(3 * time.Second)

		value, err := redisClient.Get(ctx, "ctestkey")
		assert.NoError(t, err)
		assert.Equal(t, "5838", value)

		res, err = redisClient.SetNX(ctx, "ctestkey", 5838, 5*time.Second)
		assert.NoError(t, err)
		assert.False(t, res)

		time.Sleep(3 * time.Second)

		exists, err := redisClient.Exists(ctx, "ctestkey")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TestRedisIntervalExpiry", func(t *testing.T) {
		_, err := redisClient.Delete(ctx, "dtestkey")
		assert.NoError(t, err)

		res, err := redisClient.SetNX(ctx, "dtestkey", 5838, 5*time.Second)
		assert.NoError(t, err)
		assert.True(t, res)

		time.Sleep(6 * time.Second)

		value, err := redisClient.Get(ctx, "dtestkey")
		assert.Error(t, err)
		assert.Equal(t, "", value)

		t.Logf("Redis value has expired")
	})
}
