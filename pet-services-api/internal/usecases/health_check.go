package usecases

import (
	"context"
	"pet-services-api/internal/logging"
)

type HealthCheckOutput struct {
	Status string `json:"status"`
}

type HealthCheckUseCase struct {
	logger logging.LoggerInterface
}

func NewHealthCheckUseCase(logger logging.LoggerInterface) *HealthCheckUseCase {
	return &HealthCheckUseCase{logger: logger}
}

func (h *HealthCheckUseCase) Execute(ctx context.Context) HealthCheckOutput {
	h.logger.LogInfo(ctx, "HealthCheckUseCase.Execute", "Health check executed successfully")

	return HealthCheckOutput{Status: "ok"}
}
