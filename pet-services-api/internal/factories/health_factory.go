package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"
)

type HealthFactory struct {
	HealthCheckAPI *usecases.HealthCheckAPIUseCase
	Logger         logging.LoggerInterface
}

func NewHealthFactory(logger logging.LoggerInterface) *HealthFactory {
	healthCheck := usecases.NewHealthCheckAPIUseCase(logger)

	return &HealthFactory{
		HealthCheckAPI: healthCheck,
		Logger:         logger,
	}
}
