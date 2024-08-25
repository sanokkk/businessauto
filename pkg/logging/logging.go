package logging

import (
	"autoshop/internal/config"
	"log/slog"
	"os"
)

const (
	dev  = "dev"
	prod = "prod"
)

func GetLogger(envConfig *config.EnvConfig) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(envConfig.Env)})

	return slog.New(handler)
}

func getLogLevel(env string) slog.Level {
	switch env {
	case dev:
		return slog.LevelDebug
	case prod:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

func CreateLoggerWithOp(operation string) *slog.Logger {
	cfg := config.MustLoadConfig()

	logger := GetLogger(&cfg.EnvConfig)
	return logger.With(slog.String("op", operation))
}
