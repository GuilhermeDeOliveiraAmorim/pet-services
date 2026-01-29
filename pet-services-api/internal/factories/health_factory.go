package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type HealthFactory struct {
	HealthCheckAPI *usecases.HealthCheckAPIUseCase
	HealthCheckDB  *usecases.HealthCheckDBUseCase
	Logger         logging.LoggerInterface
}

func NewHealthFactory(db *gorm.DB, logger logging.LoggerInterface) *HealthFactory {
	healthCheckAPI := usecases.NewHealthCheckAPIUseCase(logger)
	healthCheckDB := usecases.NewHealthCheckDBUseCase(logger, db)

	return &HealthFactory{
		HealthCheckAPI: healthCheckAPI,
		HealthCheckDB:  healthCheckDB,
		Logger:         logger,
	}
}
