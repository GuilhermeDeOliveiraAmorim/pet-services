package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReferenceHandler struct {
	ListCountriesUseCase *usecases.ListCountriesUseCase
	ListStatesUseCase    *usecases.ListStatesUseCase
	ListCitiesUseCase    *usecases.ListCitiesUseCase
	Logger               logging.LoggerInterface
}

func NewReferenceHandler(factory *factories.ReferenceFactory, logger logging.LoggerInterface) *ReferenceHandler {
	return &ReferenceHandler{
		ListCountriesUseCase: factory.ListCountries,
		ListStatesUseCase:    factory.ListStates,
		ListCitiesUseCase:    factory.ListCities,
		Logger:               logger,
	}
}

// ListCountries godoc
// @Summary Lista países
// @Tags Referências
// @Accept json
// @Produce json
// @Success 200 {object} usecases.ListCountriesOutput
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /reference/countries [get]
func (h *ReferenceHandler) ListCountries(c *gin.Context) {
	ctx := c.Request.Context()
	output, errs := h.ListCountriesUseCase.Execute(ctx)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ListStates godoc
// @Summary Lista estados
// @Tags Referências
// @Accept json
// @Produce json
// @Success 200 {object} usecases.ListStatesOutput
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /reference/states [get]
func (h *ReferenceHandler) ListStates(c *gin.Context) {
	ctx := c.Request.Context()
	output, errs := h.ListStatesUseCase.Execute(ctx)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ListCities godoc
// @Summary Lista cidades
// @Tags Referências
// @Accept json
// @Produce json
// @Param state_id query int false "ID do estado (IBGE)"
// @Success 200 {object} usecases.ListCitiesOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /reference/cities [get]
func (h *ReferenceHandler) ListCities(c *gin.Context) {
	ctx := c.Request.Context()
	stateIDQuery := c.Query("state_id")

	var stateID *int
	if stateIDQuery != "" {
		parsed, err := strconv.Atoi(stateIDQuery)
		if err != nil {
			problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "state_id inválido",
				Detail: "state_id deve ser numérico",
			})
			h.Logger.LogBadRequest(ctx, "ReferenceHandler.ListCities", problem.Detail, err)
			c.AbortWithStatusJSON(http.StatusBadRequest, problem)
			return
		}
		stateID = &parsed
	}

	output, errs := h.ListCitiesUseCase.Execute(ctx, usecases.ListCitiesInput{StateID: stateID})
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}
