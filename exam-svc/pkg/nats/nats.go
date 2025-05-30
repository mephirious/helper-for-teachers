package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn *nats.Conn
}

func NewClient(natsURL string) (*Client, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conn: conn,
	}

	return client, nil
}

func Connect(url string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	client := &Client{
		Conn: nc,
	}
	return client, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
