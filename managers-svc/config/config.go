package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		Server   Server
		NATS     NATSConfig
		Database DatabaseConfig
		SMTP     SMTPConfig
	}

	Server struct {
		Address string
		Port    string
	}

	NATSConfig struct {
		URL       string
		ClusterID string
		ClientID  string
	}

	DatabaseConfig struct {
		DBUser     string
		DBPassword string
		DBHost     string
		DBPort     string
		DBName     string
	}

	GRPCConfig struct {
		Port string
	}

	SMTPConfig struct {
		Host     string
		Port     int64
		Email    string
		Password string
	}
)

func NewConfig(filename string) *Config {
	if filename != "" {
		if err := loadEnv(filename); err != nil {
			panic(fmt.Sprintf("Error loading .env file: %v", err))
		}
	}

	return &Config{
		Server{
			Address: getEnvStr("ADDRESS", "localhost"),
			Port:    getEnvStr("GRPC_PORT", "50055"),
		},
		NATSConfig{
			URL:       getEnvStr("NATS_URL", "nats://localhost:4222"),
			ClusterID: getEnvStr("NATS_CLUSTER_ID", "test-cluster"),
			ClientID:  getEnvStr("NATS_CLIENT_ID", "statistics-service"),
		},
		DatabaseConfig{
			DBUser:     getEnvStr("POSTGRES_USER", "admin"),
			DBPassword: getEnvStr("POSTGRES_PASSWORD", "adminadmin"),
			DBHost:     getEnvStr("DB_HOST", "localhost"),
			DBPort:     getEnvStr("DB_PORT", "5434"),
			DBName:     getEnvStr("POSTGRES_DB", "stats_db"),
		},
		SMTPConfig{
			Host:     getEnvStr("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt64("SMTP_POST", 587),
			Email:    getEnvStr("SMTP_EMAIL", "example@gmail.com"),
			Password: getEnvStr("SMTP_PASSWORD", "YOUR_PASSWORD"),
		},
	}
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func (d *DatabaseConfig) MakeConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.DBHost, d.DBPort, d.DBUser, d.DBPassword, d.DBName,
	)
}

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
