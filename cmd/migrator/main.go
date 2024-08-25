package main

import (
	"autoshop/internal/storage/migrate"
	"flag"
)

func main() {
	var direction int
	flag.IntVar(&direction, "direction", 1, "Use 1 to forward, 2 to rollback last migration")

	flag.Parse()
	migrate.Migrate(direction)
}
