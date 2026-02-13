package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type RequestFactory struct {
	AddRequest      *usecases.AddRequestUseCase
	ListRequests    *usecases.ListRequestsUseCase
	AcceptRequest   *usecases.AcceptRequestUseCase
	CompleteRequest *usecases.CompleteRequestUseCase
}

func NewRequestFactory(db *gorm.DB, storageService storage.ObjectStorage, logger logging.LoggerInterface) *RequestFactory {
	userRepo := repository_impl.NewUserRepository(db)
	petRepo := repository_impl.NewPetRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	requestRepo := repository_impl.NewRequestRepository(db)

	return &RequestFactory{
		AddRequest:      usecases.NewAddRequestUseCase(userRepo, petRepo, serviceRepo, providerRepo, requestRepo, logger),
		ListRequests:    usecases.NewListRequestsUseCase(userRepo, requestRepo, providerRepo, serviceRepo, storageService, logger),
		AcceptRequest:   usecases.NewAcceptRequestUseCase(userRepo, providerRepo, requestRepo, logger),
		CompleteRequest: usecases.NewCompleteRequestUseCase(userRepo, providerRepo, requestRepo, logger),
	}
}
