package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"pet-services-api/internal/database"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, sqlDB := database.SetupDatabaseConnection(ctx)
	if db == nil {
		fmt.Fprintln(os.Stderr, "warning: database connection not available, exiting")
		os.Exit(1)
	}
	defer database.Shutdown(ctx, db)
	defer sqlDB.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("[Start] failed to run database migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("[Start] database migrations completed successfully")

	slog.Info("[Start] service is running. Press Ctrl+C to stop.")

	<-ctx.Done()
	slog.Info("[Start] shutting down gracefully...")
}
