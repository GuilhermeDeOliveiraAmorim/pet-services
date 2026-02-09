package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	PetFactory *factories.PetFactory
	Logger     logging.LoggerInterface
}

func NewPetHandler(factory *factories.PetFactory, logger logging.LoggerInterface) *PetHandler {
	return &PetHandler{
		PetFactory: factory,
		Logger:     logger,
	}
}

// AddPet godoc
// @Summary Adiciona um novo pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param input body usecases.AddPetInputBody true "Dados do pet"
// @Success 201 {object} usecases.AddPetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets [post]
func (h *PetHandler) AddPet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddPetInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do pet",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPet", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddPetInput{
		UserID:   userID.(string),
		Name:     inputBody.Name,
		SpecieID: inputBody.SpecieID,
		Age:      inputBody.Age,
		Weight:   inputBody.Weight,
		Notes:    inputBody.Notes,
	}

	output, errs := h.PetFactory.AddPet.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
