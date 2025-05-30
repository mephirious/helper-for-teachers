package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	// Start an in-memory Redis server
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")
	defer mr.Close()

	ctx := context.Background()
	cfg := Config{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
		TTL:      24 * time.Hour,
	}

	t.Run("successful connection", func(t *testing.T) {
		client, err := NewClient(ctx, cfg)
		require.NoError(t, err, "expected no error when connecting to Redis")
		assert.NotNil(t, client, "client should not be nil")
		assert.NotNil(t, client.client, "client.client should not be nil")
		assert.Equal(t, cfg.TTL, client.ttl, "TTL should match config")
	})

	t.Run("invalid address", func(t *testing.T) {
		cfg.Addr = "invalid:6379"
		client, err := NewClient(ctx, cfg)
		assert.Error(t, err, "expected error for invalid Redis address")
		assert.Contains(t, err.Error(), "failed to ping Redis", "error should mention ping failure")
		assert.Nil(t, client, "client should be nil on error")
	})
}

func TestClient_Close(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")
	defer mr.Close()

	ctx := context.Background()
	cfg := Config{Addr: mr.Addr(), TTL: 24 * time.Hour}
	client, err := NewClient(ctx, cfg)
	require.NoError(t, err, "failed to create client")

	err = client.Close()
	assert.NoError(t, err, "expected no error when closing client")
}

func TestClient_Ping(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")
	defer mr.Close()

	ctx := context.Background()
	cfg := Config{Addr: mr.Addr(), TTL: 24 * time.Hour}
	client, err := NewClient(ctx, cfg)
	require.NoError(t, err, "failed to create client")

	t.Run("successful ping", func(t *testing.T) {
		err := client.Ping(ctx)
		assert.NoError(t, err, "expected no error when pinging Redis")
	})

	t.Run("ping after close", func(t *testing.T) {
		client.Close()
		err := client.Ping(ctx)
		assert.Error(t, err, "expected error when pinging closed client")
		assert.Contains(t, err.Error(), "redis ping error", "error should mention ping failure")
	})
}

func TestClient_Unwrap(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")
	defer mr.Close()

	ctx := context.Background()
	cfg := Config{Addr: mr.Addr(), TTL: 24 * time.Hour}
	client, err := NewClient(ctx, cfg)
	require.NoError(t, err, "failed to create client")

	redisClient := client.Unwrap()
	assert.NotNil(t, redisClient, "Unwrap should return non-nil redis.Client")
	assert.IsType(t, &redis.Client{}, redisClient, "Unwrap should return *redis.Client")
}

func TestClient_TTL(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")
	defer mr.Close()

	ctx := context.Background()
	ttl := 48 * time.Hour
	cfg := Config{Addr: mr.Addr(), TTL: ttl}
	client, err := NewClient(ctx, cfg)
	require.NoError(t, err, "failed to create client")

	assert.Equal(t, ttl, client.TTL(), "TTL should match configured value")
}
