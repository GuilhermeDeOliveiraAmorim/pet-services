package logging

import (
	"context"
	"log/slog"
	"time"
)

// EnsureLogger garante um logger não-nil.
func EnsureLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}

// UseCase registra início e fim de um caso de uso.
func UseCase(ctx context.Context, logger *slog.Logger, usecase string, attrs ...slog.Attr) func(err *error) {
	start := time.Now()
	l := EnsureLogger(logger).With("usecase", usecase)
	if len(attrs) > 0 {
		args := make([]any, 0, len(attrs))
		for _, a := range attrs {
			args = append(args, a)
		}
		l = l.With(args...)
	}
	l.InfoContext(ctx, "start")

	return func(err *error) {
		duration := time.Since(start)
		if err != nil && *err != nil {
			l.ErrorContext(ctx, "end", slog.String("status", "error"), slog.Duration("duration", duration), slog.Any("error", *err))
			return
		}
		l.InfoContext(ctx, "end", slog.String("status", "ok"), slog.Duration("duration", duration))
	}
}
