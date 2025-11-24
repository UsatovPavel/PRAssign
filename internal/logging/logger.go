package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	hdl := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(hdl)
}
