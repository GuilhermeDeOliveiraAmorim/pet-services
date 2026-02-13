package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type ProviderFactory struct {
	AddProvider      *usecases.AddProviderUseCase
	AddProviderPhoto *usecases.AddProviderPhotoUseCase
	GetProvider      *usecases.GetProviderUseCase
}

func NewProviderFactory(db *gorm.DB, storageService storage.ObjectStorage, logger logging.LoggerInterface) *ProviderFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	photoRepo := repository_impl.NewPhotoRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)

	return &ProviderFactory{
		AddProvider:      usecases.NewAddProviderUseCase(userRepo, providerRepo, logger),
		AddProviderPhoto: usecases.NewAddProviderPhotoUseCase(userRepo, providerRepo, photoRepo, storageService, logger),
		GetProvider:      usecases.NewGetProviderUseCase(providerRepo, serviceRepo, storageService, logger),
	}
}
