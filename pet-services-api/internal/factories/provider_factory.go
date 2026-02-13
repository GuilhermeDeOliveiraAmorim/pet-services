package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type ProviderFactory struct {
	AddProvider *usecases.AddProviderUseCase
}

func NewProviderFactory(db *gorm.DB, logger logging.LoggerInterface) *ProviderFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)

	return &ProviderFactory{
		AddProvider: usecases.NewAddProviderUseCase(userRepo, providerRepo, logger),
	}
}
