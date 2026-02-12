package main

import (
	"log/slog"
	"os"
)

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(slogHandler)

	slog.SetDefault(logger)
}
