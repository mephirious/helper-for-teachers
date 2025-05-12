package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/pkg/mongo"
)

type (
	Config struct {
		Mongo   mongo.Config
		NATS    NATSConfig
		Server  Server
		Mailjet MailjetConfig
	}

	Server struct {
		GRPCServer GRPCServer
	}

	GRPCServer struct {
		Port int
		Mode string
	}

	NATSConfig struct {
		URL string
	}

	MailjetConfig struct {
		API  string
		KEY  string
		From string
		Name string
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: error loading .env file: %v", err)
	}

	var cfg Config

	port := getEnv("GRPC_PORT", "8080")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT value: %w", err)
	}
	cfg.Server.GRPCServer.Port = portInt
	cfg.Server.GRPCServer.Mode = getEnv("GIN_MODE", "release")

	cfg.Mongo.Database = os.Getenv("MONGO_DB")
	cfg.Mongo.URI = os.Getenv("MONGO_DB_URI")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")

	cfg.NATS.URL = os.Getenv("NATS_URL")

	cfg.Mailjet.API = os.Getenv("MAILJET_API_KEY")
	cfg.Mailjet.KEY = os.Getenv("MAILJET_SECRET_KEY")
	cfg.Mailjet.From = os.Getenv("MAILJET_FROM_EMAIL")
	cfg.Mailjet.Name = os.Getenv("MAILJET_FROM_NAME")

	return &cfg, nil
}

func getEnv(field, defaultVal string) string {
	value := os.Getenv(field)
	if value == "" {
		return defaultVal
	}
	return value
}
