package usecases

import (
	"context"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type HealthCheckDBOutput struct {
	Status string `json:"status"`
}

type HealthCheckDBUseCase struct {
	logger logging.LoggerInterface
	db     *gorm.DB
}

func NewHealthCheckDBUseCase(logger logging.LoggerInterface, db *gorm.DB) *HealthCheckDBUseCase {
	return &HealthCheckDBUseCase{logger: logger, db: db}
}

func (h *HealthCheckDBUseCase) Execute(ctx context.Context) HealthCheckDBOutput {
	var status string

	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.LogError(ctx, "HealthCheckDBUseCase.Execute", "Erro ao obter conexão DB", err)
		status = "error"
	} else if err = sqlDB.Ping(); err != nil {
		h.logger.LogError(ctx, "HealthCheckDBUseCase.Execute", "Ping ao banco falhou", err)
		status = "error"
	} else {
		h.logger.LogInfo(ctx, "HealthCheckDBUseCase.Execute", "Health check for DB executado com sucesso")
		status = "ok"
	}
	return HealthCheckDBOutput{Status: status}
}
