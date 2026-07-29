package logger

import (
	"log/slog"
	"os"
)

func New(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func NewDefault() *slog.Logger {
	return New(slog.LevelInfo)
}
