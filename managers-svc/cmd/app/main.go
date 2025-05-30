package main

import (
	"os"

	"github.com/mephirious/helper-for-teachers/managers-svc/config"
	"github.com/mephirious/helper-for-teachers/managers-svc/internal/app"
	"github.com/mephirious/helper-for-teachers/managers-svc/pkg/lib/prettyslog"
)

func main() {
	cfg := config.NewConfig(".env")
	logger := prettyslog.SetupPrettySlog(os.Stdout)

	server := app.NewAPIServer(cfg, logger)
	server.Run()
}
