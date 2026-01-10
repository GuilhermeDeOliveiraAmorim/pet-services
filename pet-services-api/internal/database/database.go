package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	POSTGRES = "postgres"
)

type StorageInput struct {
	DB         *gorm.DB
	BucketName string
}

type gormLogAdapter struct{}

func (g gormLogAdapter) Printf(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

func NewDatabase(ctx context.Context) *gorm.DB {
	newLogger := logger.New(
		gormLogAdapter{},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" {
		slog.Error("one or more required database environment variables are not set")
		return nil
	}

	// Validate DB_PORT format if provided
	if dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err != nil || port < 1 || port > 65535 {
			slog.Error("invalid DB_PORT: must be a number between 1 and 65535", "value", dbPort)
			return nil
		}
	}

	var dsn string
	if dbPort == "" {
		dsn = "postgresql://" + dbUser + ":" + dbPass + "@" + dbHost + "/" + dbName + "?sslmode=require"
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
			dbHost, dbUser, dbPass, dbName, dbPort)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil
	}

	return db
}

func SetupDatabaseConnection(ctx context.Context) (*gorm.DB, *sql.DB) {
	db := NewDatabase(ctx)
	if db == nil {
		slog.Error("database connection not available")
		return nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm.DB", "error", err)
		return nil, nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(2 * time.Minute)

	return db, sqlDB
}

func CheckConnection(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm.DB", "error", err)
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		slog.Error("database ping failed", "error", err)
		return false
	}

	return true
}

func Shutdown(ctx context.Context, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm.DB", "error", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
	}
}
