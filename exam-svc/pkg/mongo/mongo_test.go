package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := Config{
		URI:      "mongodb://localhost:27017",
		Database: "testdb",
		Username: "",
		Password: "",
	}

	db, err := NewDB(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NotNil(t, db.Connection)
	require.NotNil(t, db.Client)

	err = db.Client.Disconnect(ctx)
	require.NoError(t, err)
}
