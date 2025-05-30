package nats

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startTestNATSServer(t *testing.T) *server.Server {
	t.Helper()
	opts := &server.Options{
		Host: "localhost",
		Port: -1,
	}
	natsServer, err := server.NewServer(opts)
	require.NoError(t, err, "failed to create test NATS server")
	natsServer.Start()
	if !natsServer.ReadyForConnections(5 * time.Second) {
		t.Fatal("NATS server failed to start")
	}
	t.Cleanup(func() {
		natsServer.Shutdown()
	})
	return natsServer
}

func TestNewClient(t *testing.T) {
	natsServer := startTestNATSServer(t)
	url := natsServer.ClientURL()

	t.Run("successful connection", func(t *testing.T) {
		client, err := NewClient(url)
		require.NoError(t, err, "expected no error when connecting to NATS")
		assert.NotNil(t, client, "client should not be nil")
		assert.NotNil(t, client.Conn, "client.Conn should not be nil")
		assert.True(t, client.Conn.IsConnected(), "client should be connected")
		client.Close()
	})

	t.Run("invalid URL", func(t *testing.T) {
		client, err := NewClient("invalid://url")
		assert.Error(t, err, "expected error for invalid NATS URL")
		assert.Nil(t, client, "client should be nil on error")
	})
}

func TestConnect(t *testing.T) {
	natsServer := startTestNATSServer(t)
	url := natsServer.ClientURL()

	t.Run("successful connection", func(t *testing.T) {
		client, err := Connect(url)
		require.NoError(t, err, "expected no error when connecting to NATS")
		assert.NotNil(t, client, "client should not be nil")
		assert.NotNil(t, client.Conn, "client.Conn should not be nil")
		assert.True(t, client.Conn.IsConnected(), "client should be connected")
		client.Close()
	})

	t.Run("invalid URL", func(t *testing.T) {
		client, err := Connect("invalid://url")
		assert.Error(t, err, "expected error for invalid NATS URL")
		assert.Contains(t, err.Error(), "failed to connect to NATS", "error should mention connection failure")
		assert.Nil(t, client, "client should be nil on error")
	})
}

func TestClient_Close(t *testing.T) {
	natsServer := startTestNATSServer(t)
	url := natsServer.ClientURL()

	client, err := NewClient(url)
	require.NoError(t, err, "failed to create client")

	t.Run("close connected client", func(t *testing.T) {
		assert.True(t, client.Conn.IsConnected(), "client should be connected initially")
		client.Close()
		assert.False(t, client.Conn.IsConnected(), "client should be disconnected after Close")
	})

	t.Run("close nil connection", func(t *testing.T) {
		c := &Client{Conn: nil}
		c.Close()
	})
}
