package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	migrations "pet-services-api/internal/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds connection and migration options.
type Config struct {
	DSN             string
	RunMigrations   bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Logger          *slog.Logger
	LogLevel        logger.LogLevel
}

// DefaultConfig builds a config from environment variables.
func DefaultConfig() Config {
	return Config{
		DSN:             os.Getenv("DATABASE_URL"),
		RunMigrations:   true,
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 2 * time.Minute,
		LogLevel:        logger.Silent,
	}
}

// Open connects to Postgres via GORM, pings, and optionally runs migrations.
func Open(ctx context.Context, cfg Config) (*gorm.DB, *sql.DB, error) {
	dsn, err := resolveDSN(cfg)
	if err != nil {
		return nil, nil, err
	}

	slogger := cfg.Logger
	if slogger == nil {
		slogger = slog.Default()
	}

	gormLogger := logger.New(
		slogAdapter{logger: slogger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  cfg.LogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("database pool handle: %w", err)
	}

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}

	if cfg.RunMigrations {
		if err := migrations.RunMigrations(db.WithContext(ctx)); err != nil {
			_ = sqlDB.Close()
			return nil, nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	return db, sqlDB, nil
}

// Close attempts to close the underlying sql.DB.
func Close(logger *slog.Logger, sqlDB *sql.DB) {
	if sqlDB == nil {
		return
	}
	if err := sqlDB.Close(); err != nil {
		if logger == nil {
			logger = slog.Default()
		}
		logger.Error("failed to close database connection", "error", err)
	}
}

// WithTx runs fn inside a transaction using the provided context.
func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if db == nil {
		return errors.New("db is nil")
	}
	return db.WithContext(ctx).Transaction(fn)
}

func resolveDSN(cfg Config) (string, error) {
	if cfg.DSN != "" {
		return cfg.DSN, nil
	}
	if env := os.Getenv("DATABASE_URL"); env != "" {
		return env, nil
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	slog.Info("Variáveis", user, pass, name, host, port, sslMode)

	if sslMode == "" {
		sslMode = "disable"
	}

	if user == "" || pass == "" || name == "" || host == "" {
		return "", errors.New("database dsn not provided; set DATABASE_URL or DB_USER/DB_PASS/DB_NAME/DB_HOST")
	}

	if port == "" {
		port = "5432"
	} else if _, err := strconv.Atoi(port); err != nil {
		return "", fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, pass, name, port, sslMode), nil
}

type slogAdapter struct {
	logger *slog.Logger
}

func (a slogAdapter) Printf(format string, args ...any) {
	a.logger.Info(fmt.Sprintf(format, args...))
}
