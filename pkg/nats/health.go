package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func (c *Client) HealthCheck(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Update nats state
	if err := c.conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("nats flush: %w", err)
	}
	// Check the state
	if c.conn.Status() != nats.CONNECTED {
		return fmt.Errorf("nats status: %s", c.conn.Status())
	}
	return nil
}
