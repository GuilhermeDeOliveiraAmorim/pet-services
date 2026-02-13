package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	AddService      *usecases.AddServiceUseCase
	AddServicePhoto *usecases.AddServicePhotoUseCase
}

func NewServiceFactory(db *gorm.DB, storageService storage.ObjectStorage, logger logging.LoggerInterface) *ServiceFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)
	photoRepo := repository_impl.NewPhotoRepository(db)

	return &ServiceFactory{
		AddService:      usecases.NewAddServiceUseCase(userRepo, providerRepo, serviceRepo, logger),
		AddServicePhoto: usecases.NewAddServicePhotoUseCase(userRepo, serviceRepo, providerRepo, photoRepo, storageService, logger),
	}
}
