package database

import (
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

const (
	versionInitial = "20260103000000"
	versionV2      = "20260103000001"
	versionV3      = "20260103000002"
	versionV4      = "20260103000003"
	versionV5      = "20260109000001"
	versionV6      = "20260109000002"
)

func getMigrations() []Migration {
	return []Migration{
		{
			Version:     versionInitial,
			Description: "Initial schema (users, providers, requests, reviews, tokens)",
			Up:          migrationInitialSchema,
		},
		{
			Version:     versionV2,
			Description: "Add relationships and rename service_requests to requests",
			Up:          migrationAddRelationshipsV2,
		},
		{
			Version:     versionV3,
			Description: "Add user relationships to tokens",
			Up:          migrationAddTokenRelationshipsV3,
		},
		{
			Version:     versionV4,
			Description: "Add performance indexes for faster queries",
			Up:          migrationAddIndexesV4,
		},
		{
			Version:     versionV5,
			Description: "Add provider_service_photos table for multiple service images",
			Up:          migrationAddProviderServicePhotosV5,
		},
		{
			Version:     versionV6,
			Description: "Add foreign key and relation between provider_services and provider_service_photos",
			Up:          migrationAddServicePhotoRelationV6,
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

// migrationAddRelationshipsV2 adiciona relações e renomeia service_requests para requests.
func migrationAddRelationshipsV2(db *gorm.DB) error {
	// Renomear tabela service_requests para requests
	if db.Migrator().HasTable("service_requests") {
		if err := db.Migrator().RenameTable("service_requests", "requests"); err != nil {
			return err
		}
	}

	// Reexecuta AutoMigrate para adicionar as relações definidas nos modelos
	return db.AutoMigrate(
		&models.Provider{},
		&models.ServiceRequest{},
		&models.Review{},
	)
}

// migrationAddTokenRelationshipsV3 adiciona relações de user_id aos tokens.
func migrationAddTokenRelationshipsV3(db *gorm.DB) error {
	// Reexecuta AutoMigrate para adicionar as relações definidas nos modelos
	return db.AutoMigrate(
		&models.RefreshToken{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
	)
}

// migrationAddIndexesV4 adiciona índices para melhorar a performance de consultas.
func migrationAddIndexesV4(db *gorm.DB) error {
	// Reexecuta AutoMigrate para adicionar todos os índices definidos nos modelos
	return db.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.ServiceRequest{},
		&models.Review{},
		&models.RefreshToken{},
		&models.EmailVerificationToken{},
		&models.PasswordResetToken{},
	)
}

// migrationAddProviderServicePhotosV5 cria a tabela provider_service_photos.
func migrationAddProviderServicePhotosV5(db *gorm.DB) error {
	return db.AutoMigrate(&models.ProviderServicePhoto{})
}

// migrationAddServicePhotoRelationV6 garante a relação entre provider_services e provider_service_photos.
func migrationAddServicePhotoRelationV6(db *gorm.DB) error {
	return db.AutoMigrate(&models.ProviderServicePhoto{}, &models.ProviderService{})
}
