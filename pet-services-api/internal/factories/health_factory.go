package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"
)

type HealthFactory struct {
	HealthCheck *usecases.HealthCheckUseCase
	Logger      logging.LoggerInterface
}

func NewHealthFactory(logger logging.LoggerInterface) *HealthFactory {
	healthCheck := usecases.NewHealthCheckUseCase(logger)

	return &HealthFactory{
		HealthCheck: healthCheck,
		Logger:      logger,
	}
}
