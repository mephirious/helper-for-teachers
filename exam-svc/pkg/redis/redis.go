package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string        `env:"HOST" envDefault:"localhost:6379"`
	Password string        `env:"PASSWORD"`
	DB       int           `env:"DB" envDefault:"0"`
	TTL      time.Duration `env:"TTL" envDefault:"86400s"`
}

type Client struct {
	client *redis.Client
	ttl    time.Duration
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &Client{
		client: rdb,
		ttl:    cfg.TTL,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping error: %w", err)
	}
	return nil
}

func (c *Client) Unwrap() *redis.Client {
	return c.client
}

func (c *Client) TTL() time.Duration {
	return c.ttl
}
