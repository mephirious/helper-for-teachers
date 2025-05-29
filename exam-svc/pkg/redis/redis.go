package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string        `env:"HOST"`
	Password string        `env:"PASSWORD"`
	TTL      time.Duration `env:"TTL" envDefault:"0"`
}

type Client struct {
	client *redis.Client
	ttl    time.Duration
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{
		client: client,
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
