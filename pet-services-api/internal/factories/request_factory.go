package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type RequestFactory struct {
	AddRequest *usecases.AddRequestUseCase
}

func NewRequestFactory(db *gorm.DB, logger logging.LoggerInterface) *RequestFactory {
	userRepo := repository_impl.NewUserRepository(db)
	petRepo := repository_impl.NewPetRepository(db)
	serviceRepo := repository_impl.NewServiceRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	requestRepo := repository_impl.NewRequestRepository(db)

	return &RequestFactory{
		AddRequest: usecases.NewAddRequestUseCase(userRepo, petRepo, serviceRepo, providerRepo, requestRepo, logger),
	}
}
