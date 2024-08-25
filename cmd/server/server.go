package main

import (
	"autoshop/internal/app/api"
	"autoshop/internal/config"
	"autoshop/internal/storage/migrate"
	"autoshop/pkg/logging"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

func main() {
	cfg := config.MustLoadConfig()
	//fmt.Printf("%+v\n", cfg)

	logger := logging.GetLogger(&cfg.EnvConfig)
	log := logger.With(slog.String("op", "main"))

	migrate.Migrate(1)

	log.Info("starting controllers")

	server := api.NewServer(cfg, logger)
	server.Handler.Start(cfg.ApiConfig)
}
