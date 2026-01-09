package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pet-services-api/internal/database"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL não definido")
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("falha ao conectar no banco", "error", err)
		os.Exit(1)
	}

	if err := database.RunMigrations(db.WithContext(context.Background())); err != nil {
		slog.Error("migrações falharam", "error", err)
		os.Exit(1)
	}

	slog.Info("migrações concluídas com sucesso")
}
