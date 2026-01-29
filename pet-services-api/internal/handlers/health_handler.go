package handlers

import (
	"net/http"

	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	HealthCheckAPIUseCase *usecases.HealthCheckAPIUseCase
	HealthCheckDBUseCase  *usecases.HealthCheckDBUseCase
	Logger                logging.LoggerInterface
}

func NewHealthHandler(factory *factories.HealthFactory, logger logging.LoggerInterface) *HealthHandler {
	return &HealthHandler{
		HealthCheckAPIUseCase: factory.HealthCheckAPI,
		HealthCheckDBUseCase:  factory.HealthCheckDB,
		Logger:                logger,
	}
}

type HealthCheckResponse struct {
	API usecases.HealthCheckAPIOutput `json:"api"`
	DB  usecases.HealthCheckDBOutput  `json:"db"`
}

// HealthCheck godoc
// @Summary Verifica a saúde da API
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	resultAPI := h.HealthCheckAPIUseCase.Execute(ctx)
	resultDB := h.HealthCheckDBUseCase.Execute(ctx)

	c.JSON(http.StatusOK, HealthCheckResponse{
		API: resultAPI,
		DB:  resultDB,
	})
}
