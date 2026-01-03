package database

import (
	"github.com/guilherme/pet-services-api/internal/models"
	"gorm.io/gorm"
)

const (
	versionInitial = "20260103000000"
)

func getMigrations() []Migration {
	return []Migration{
		{
			Version:     versionInitial,
			Description: "Initial schema (users, providers, requests, reviews, tokens)",
			Up:          migrationInitialSchema,
		},
	}
}

// migrationInitialSchema cria o schema base com AutoMigrate.
func migrationInitialSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.ProviderService{},
		&models.ProviderPhoto{},
		&models.ProviderWorkingHour{},
		&models.ServiceRequest{},
		&models.Review{},
		&models.RefreshToken{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
	)
}
