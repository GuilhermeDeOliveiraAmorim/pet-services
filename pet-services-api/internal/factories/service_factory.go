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
	DeleteServicePhoto *usecases.DeleteServicePhotoUseCase
	AddServiceTag      *usecases.AddServiceTagUseCase
	ListTags           *usecases.ListTagsUseCase
	AddServiceCategory *usecases.AddServiceCategoryUseCase
	ListServices       *usecases.ListServicesUseCase
	SearchServices     *usecases.SearchServicesUseCase
	GetService         *usecases.GetServiceUseCase
	UpdateService      *usecases.UpdateServiceUseCase
	DeleteService      *usecases.DeleteServiceUseCase
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
		DeleteServicePhoto: usecases.NewDeleteServicePhotoUseCase(userRepo, serviceRepo, providerRepo, photoRepo, storageService, logger),
		AddServiceTag:      usecases.NewAddServiceTagUseCase(userRepo, serviceRepo, providerRepo, tagRepo, logger),
		ListTags:           usecases.NewListTagsUseCase(tagRepo),
		AddServiceCategory: usecases.NewAddServiceCategoryUseCase(userRepo, serviceRepo, providerRepo, categoryRepo, logger),
		ListServices:       usecases.NewListServicesUseCase(serviceRepo, providerRepo, storageService, logger),
		SearchServices:     usecases.NewSearchServicesUseCase(serviceRepo, providerRepo, storageService, logger),
		GetService:         usecases.NewGetServiceUseCase(serviceRepo, storageService, logger),
		UpdateService:      usecases.NewUpdateServiceUseCase(userRepo, providerRepo, serviceRepo, storageService, logger),
		DeleteService:      usecases.NewDeleteServiceUseCase(userRepo, providerRepo, serviceRepo, logger),
	}
}
