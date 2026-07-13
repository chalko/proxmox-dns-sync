package logger

import (
	"log/slog"
	"os"
)

// InitLogger configures and sets the global slog JSON logger.
func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}
