package database

import (
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

func Migration20260110000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Species{},
		&models.Category{},
		&models.Tag{},
		&models.Photo{},

		&models.Provider{},
		&models.Pet{},
		&models.Service{},
		&models.Review{},
		&models.Request{},
	)
}

func Migration20260215000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.RefreshToken{},
	)
}

func Migration20260204000000(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.PasswordResetToken{},
	)
}
