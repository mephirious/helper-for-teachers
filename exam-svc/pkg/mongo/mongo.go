package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI      string `env:"DB_URI"`
	Database string `env:"DB"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

var clientOptions *options.ClientOptions

type DB struct {
	Connection *mongo.Database
	Client     *mongo.Client
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	clientOptions = options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connection to mongoDB Error: %w", err)
	}

	db := &DB{
		Connection: client.Database(cfg.Database),
		Client:     client,
	}

	err = db.Client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping connection mongoDB Error: %w", err)
	}

	go db.reconnectOnFailure(ctx)

	return db, nil
}

func (db *DB) reconnectOnFailure(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			err := db.Client.Ping(ctx, nil)
			if err != nil {
				log.Printf("lost connection to mongoDB: %v", err)
				db.Client, _ = mongo.Connect(ctx, clientOptions)

				err = db.Client.Ping(ctx, nil)
				if err == nil {
					log.Printf("ping to mongoDB is successful: %v", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			err := db.Client.Disconnect(ctx)
			if err != nil {
				log.Printf("mongoDB close connection error: %v", err)
				return
			}

			log.Printf("mongo connection is closed successfully")
		}
	}
}
