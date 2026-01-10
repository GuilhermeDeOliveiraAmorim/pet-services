package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type SchemaMigration struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"not null"`
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("error creating migration table: %w", err)
	}

	migrations := getMigrations()

	for _, migration := range migrations {
		var existing SchemaMigration
		result := db.Where("version = ?", migration.Version).First(&existing)

		if result.Error == nil {
			continue
		}

		if result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("error checking migration %s: %w", migration.Version, result.Error)
		}

		slog.Info("[Migration] applying migration", "version", migration.Version, "description", migration.Description)

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("error executing migration: %w", err)
			}

			record := SchemaMigration{
				Version:   migration.Version,
				AppliedAt: time.Now(),
			}
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("error registering migration: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("migration failed %s: %w", migration.Version, err)
		}

		slog.Info("[Migration] migration applied successfully", "version", migration.Version)
	}

	slog.Info("[Migration] all migrations have been applied")
	return nil
}

type Migration struct {
	Version     string
	Description string
	Up          func(*gorm.DB) error
}

func getMigrations() []Migration {
	return []Migration{
		{
			Version:     "20260110000000",
			Description: "Create initial schema for pet-services",
			Up:          Migration20260110000000,
		},
	}
}
