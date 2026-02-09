package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type SpecieHandler struct {
	ListSpeciesUseCase *usecases.ListSpeciesUseCase
	Logger             logging.LoggerInterface
}

func NewSpecieHandler(factory *factories.SpecieFactory, logger logging.LoggerInterface) *SpecieHandler {
	return &SpecieHandler{
		ListSpeciesUseCase: factory.ListSpecies,
		Logger:             logger,
	}
}

// ListSpecies godoc
// @Summary Lista espécies
// @Tags Espécies
// @Accept json
// @Produce json
// @Success 200 {object} usecases.ListSpeciesOutput
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /species [get]
func (h *SpecieHandler) ListSpecies(c *gin.Context) {
	ctx := c.Request.Context()
	output, errs := h.ListSpeciesUseCase.Execute(ctx)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}
