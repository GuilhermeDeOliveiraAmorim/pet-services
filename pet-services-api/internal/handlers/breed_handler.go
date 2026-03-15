package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type BreedHandler struct {
	ListBreedsUseCase *usecases.ListBreedsUseCase
	Logger            logging.LoggerInterface
}

func NewBreedHandler(factory *factories.BreedFactory, logger logging.LoggerInterface) *BreedHandler {
	return &BreedHandler{
		ListBreedsUseCase: factory.ListBreeds,
		Logger:            logger,
	}
}

// ListBreeds godoc
// @Summary Lista raças por espécie
// @Tags Raças
// @Accept json
// @Produce json
// @Param species_id path string true "ID da espécie"
// @Success 200 {object} usecases.ListBreedsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /species/{species_id}/breeds [get]
func (h *BreedHandler) ListBreeds(c *gin.Context) {
	ctx := c.Request.Context()

	input := usecases.ListBreedsInput{SpeciesID: c.Param("species_id")}
	output, errs := h.ListBreedsUseCase.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
