package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/mongo"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/redis"
)

type (
	Config struct {
		Mongo   mongo.Config `envPrefix:"MONGO_"`
		NATS    NATSConfig   `envPrefix:"NATS_"`
		Server  Server
		Mailjet MailjetConfig `envPrefix:"MAILJET_"`
		Redis   redis.Config  `envPrefix:"REDIS_"`
		Gemini  GeminiConfig  `envPrefix:"GEMINI_"`
	}

	Server struct {
		GRPCServer GRPCServer
	}

	GRPCServer struct {
		Port int    `env:"GRPC_PORT" envDefault:"8080"`
		Mode string `env:"GIN_MODE" envDefault:"release"`
	}

	GeminiConfig struct {
		APIKey string `env:"API_KEY"`
	}

	NATSConfig struct {
		URL string `env:"URL"`
	}

	MailjetConfig struct {
		API  string `env:"API_KEY"`
		KEY  string `env:"SECRET_KEY"`
		From string `env:"FROM_EMAIL"`
		Name string `env:"FROM_NAME"`
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment: %w", err)
	}

	cfg.Redis.TTL = time.Duration(cfg.Redis.TTL.Seconds()) * time.Second

	return cfg, nil
}
