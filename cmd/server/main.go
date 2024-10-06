package main

import (
	_ "autoshop/docs"
	"autoshop/internal/app/api"
	"autoshop/internal/config"
	"autoshop/internal/storage/migrate"
	"autoshop/pkg/logging"
	"log/slog"
)

// @title BusinessAuto API
func main() {
	cfg := config.MustLoadConfig()

	logger := logging.GetLogger(&cfg.EnvConfig)
	log := logger.With(slog.String("op", "main"))

	migrate.Migrate(1)

	log.Info("starting controllers")

	server := api.NewServer(cfg, logger)
	server.Handler.Start(cfg.ApiConfig)
}
