package handlers

import (
	"net/http"

	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	HealthCheckUseCase *usecases.HealthCheckAPIUseCase
	Logger             logging.LoggerInterface
}

func NewHealthHandler(factory *factories.HealthFactory, logger logging.LoggerInterface) *HealthHandler {
	return &HealthHandler{HealthCheckUseCase: factory.HealthCheckAPI, Logger: logger}
}

// HealthCheckAPI godoc
// @Summary Verifica a saúde da API
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} usecases.HealthCheckAPIOutput
// @Router /health [get]
func (h *HealthHandler) HealthCheckAPI(c *gin.Context) {
	ctx := c.Request.Context()

	result := h.HealthCheckUseCase.Execute(ctx)

	c.JSON(http.StatusOK, result)
}
