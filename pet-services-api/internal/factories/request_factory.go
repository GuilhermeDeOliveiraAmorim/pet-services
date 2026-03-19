package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
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
	RejectRequest   *usecases.RejectRequestUseCase
	GetRequest      *usecases.GetRequestUseCase
}

func NewRequestFactory(db *gorm.DB, storageService storage.ObjectStorage, mailService mail.EmailService, logger logging.LoggerInterface) *RequestFactory {
	userRepo := repository_impl.NewUserRepository(db)
	petRepo := repository_impl.NewPetRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	requestRepo := repository_impl.NewRequestRepository(db)

	return &RequestFactory{
		AddRequest:      usecases.NewAddRequestUseCase(userRepo, petRepo, serviceRepo, providerRepo, requestRepo, mailService, logger),
		ListRequests:    usecases.NewListRequestsUseCase(userRepo, requestRepo, providerRepo, serviceRepo, storageService, logger),
		AcceptRequest:   usecases.NewAcceptRequestUseCase(userRepo, providerRepo, requestRepo, mailService, logger),
		CompleteRequest: usecases.NewCompleteRequestUseCase(userRepo, providerRepo, requestRepo, mailService, logger),
		RejectRequest:   usecases.NewRejectRequestUseCase(userRepo, providerRepo, requestRepo, mailService, logger),
		GetRequest:      usecases.NewGetRequestUseCase(userRepo, providerRepo, requestRepo, storageService, logger),
	}
}
