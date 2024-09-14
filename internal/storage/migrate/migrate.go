package migrate

import (
	"autoshop/internal/config"
	"autoshop/pkg/logging"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	forward  = 1
	rollback = 2
)

func Migrate(direction int) {
	const op = "migrate.Migrate"
	log := logging.CreateLoggerWithOp(op)

	switch direction {
	case forward:
		log.Info("Начинаю миграцию вперед")

		fullForward()
		return
	case rollback:
		log.Info("Начинаю миграцию назад")

		rollbackStep()
		return
	default:
		log.Error("Задано неправильное направление миграции")
		panic("Задано неправильное направление миграции")
	}
}

func fullForward() {
	const op = "migrate.FullForward"
	log := logging.CreateLoggerWithOp(op)

	dbConfig := config.MustLoadConfig().DbConfig

	m, err := migrate.New(
		"file://"+dbConfig.MigrationsPath,
		fmt.Sprintf("%s?sslmode=%s&x-migrations-table=%s", dbConfig.DbConnectionString, dbConfig.SslMode, dbConfig.MigrationsTable))

	if err != nil {
		log.Error(errors.Wrap(err, "Ошибка миграции").Error())

		panic("Ошибка миграции")
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		log.Error(errors.Wrap(err, "Ошибка миграции вперед").Error())

		panic("Ошибка миграции")
	}

	log.Info("Миграции выполнены")
}

func rollbackStep() {
	const op = "migrate.Rollback"
	log := logging.CreateLoggerWithOp(op)

	dbConfig := config.MustLoadConfig().DbConfig

	m, err := migrate.New(
		"file://"+dbConfig.MigrationsPath,
		fmt.Sprintf("%s?sslmode=%s&x-migrations-table=%s", dbConfig.DbConnectionString, dbConfig.SslMode, dbConfig.MigrationsTable))

	if err != nil {
		log.Error(errors.Wrap(err, "Ошибка миграции назад").Error())

		panic("Ошибка миграции назад")
	}

	if err := m.Steps(-1); err != nil {
		log.Error(errors.Wrap(err, "Ошибка миграции назад").Error())

		panic("Ошибка миграции назад")
	}

	log.Info("Миграция на шаг назад выполнена")
}
