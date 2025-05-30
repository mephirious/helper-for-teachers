package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WithDotEnv(t *testing.T) {
	envFile, err := os.CreateTemp("", "test.env")
	require.NoError(t, err, "failed to create temp .env file")
	defer os.Remove(envFile.Name())

	content := `
	MONGO_URI=mongodb://localhost:27017
	MONGO_DATABASE=testdb
	NATS_URL=nats://localhost:4222
	REDIS_HOST=localhost:6379
	REDIS_TTL=7200s
	GRPC_PORT=9090
	`
	_, err = envFile.Write([]byte(content))
	require.NoError(t, err, "failed to write to temp .env file")
	require.NoError(t, envFile.Close(), "failed to close temp .env file")

	os.Setenv("DOTENV_PATH", envFile.Name())

	for _, key := range []string{"MONGO_URI", "MONGO_DATABASE", "NATS_URL", "REDIS_HOST", "REDIS_TTL", "GRPC_PORT"} {
		t.Logf("%s=%s", key, os.Getenv(key))
	}

	cfg, err := New()
	require.NoError(t, err, "expected no error when loading .env file")

	assert.Equal(t, "", cfg.Mongo.URI, "Mongo URI mismatch")
	assert.Equal(t, "", cfg.Mongo.Database, "Mongo Database mismatch")
	assert.Equal(t, "", cfg.NATS.URL, "NATS URL mismatch")
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr, "Redis Addr mismatch")
	assert.Equal(t, 86400*time.Second, cfg.Redis.TTL, "Redis TTL mismatch")
	assert.Equal(t, 8080, cfg.Server.GRPCServer.Port, "GRPC Port mismatch")
}
