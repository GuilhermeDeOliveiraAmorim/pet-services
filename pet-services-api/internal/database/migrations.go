package database

import (
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

// Migration20260110000000 cria todas as tabelas do sistema pet-services
func Migration20260110000000(db *gorm.DB) error {
	return db.AutoMigrate(
		// Entidades base
		&models.User{},
		&models.Specie{},
		&models.Breed{},
		&models.Category{},
		&models.Tag{},
		&models.Photo{},

		// Entidades dependentes
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

// Migration20260124000000 adiciona o campo token_type à tabela refresh_tokens
func Migration20260124000000(db *gorm.DB) error {
	return db.Migrator().AddColumn(&models.RefreshToken{}, "token_type")
}