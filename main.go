package main

import (
	"log/slog"

	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
)

func main() {
	err := Start()
	if err != nil {
		slog.Error("Error running the server", logger.Err(err))
		panic(err)
	}
}
