package usecases

import (
	"context"
	"pet-services-api/internal/logging"
)

type HealthCheckAPIOutput struct {
	Status string `json:"status"`
}

type HealthCheckAPIUseCase struct {
	logger logging.LoggerInterface
}

func NewHealthCheckAPIUseCase(logger logging.LoggerInterface) *HealthCheckAPIUseCase {
	return &HealthCheckAPIUseCase{logger: logger}
}

func (h *HealthCheckAPIUseCase) Execute(ctx context.Context) HealthCheckAPIOutput {
	h.logger.LogInfo(ctx, "HealthCheckAPIUseCase.Execute", "Health check executed successfully")

	return HealthCheckAPIOutput{Status: "ok"}
}
