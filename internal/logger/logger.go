package logger

import (
	"log/slog"
	"os"
	)

func New() *slog.Logger {
	level := slog.LevelInfo

	if os.Getenv("ENV") == "dev" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}