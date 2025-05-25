package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mephirious/helper-for-teachers/services/auth-svc/config"
	"github.com/mephirious/helper-for-teachers/services/auth-svc/internal/app"
	"github.com/mephirious/helper-for-teachers/services/auth-svc/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		println("No .env file found, falling back to real env")
	}

	cfg, err := config.New()
	if err != nil {
		stdlog.Fatalf("config load error: %v", err)
	}

	// logger
	log := logger.New(logger.Config(cfg.Log))
	log.Info("config loaded", "version", cfg.Version)

	// ctx to cancel build or run phase, and proceed to shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Build the app
	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("app init failed", "err", err)
	}

	// Run in a goroutine so we can also listen for ctx.Done() or appErr
	appErr := make(chan error)
	go func() {
		appErr <- app.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-appErr:
		if err != nil {
			log.Error("app run error", "err", err)
		}
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Error("error during shutdown", "err", err)
	}

	log.Info("Authentification service stopped")
}
