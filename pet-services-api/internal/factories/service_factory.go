package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	AddService         *usecases.AddServiceUseCase
	AddServicePhoto    *usecases.AddServicePhotoUseCase
	AddServiceTag      *usecases.AddServiceTagUseCase
	ListTags           *usecases.ListTagsUseCase
	AddServiceCategory *usecases.AddServiceCategoryUseCase
	ListServices       *usecases.ListServicesUseCase
}

func NewServiceFactory(db *gorm.DB, storageService storage.ObjectStorage, logger logging.LoggerInterface) *ServiceFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)
	photoRepo := repository_impl.NewPhotoRepository(db)
	tagRepo := repository_impl.NewTagRepository(db)
	categoryRepo := repository_impl.NewCategoryRepository(db)

	return &ServiceFactory{
		AddService:         usecases.NewAddServiceUseCase(userRepo, providerRepo, serviceRepo, logger),
		AddServicePhoto:    usecases.NewAddServicePhotoUseCase(userRepo, serviceRepo, providerRepo, photoRepo, storageService, logger),
		AddServiceTag:      usecases.NewAddServiceTagUseCase(userRepo, serviceRepo, providerRepo, tagRepo, logger),
		ListTags:           usecases.NewListTagsUseCase(tagRepo),
		AddServiceCategory: usecases.NewAddServiceCategoryUseCase(userRepo, serviceRepo, providerRepo, categoryRepo, logger),
		ListServices:       usecases.NewListServicesUseCase(serviceRepo, providerRepo, storageService, logger),
	}
}
