package handlers

import (
	"net/http"

	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AdoptionGuardianHandler struct {
	AdoptionGuardianFactory *factories.AdoptionGuardianFactory
	Logger                  logging.LoggerInterface
}

func NewAdoptionGuardianHandler(factory *factories.AdoptionGuardianFactory, logger logging.LoggerInterface) *AdoptionGuardianHandler {
	return &AdoptionGuardianHandler{
		AdoptionGuardianFactory: factory,
		Logger:                  logger,
	}
}

// CreateAdoptionGuardianProfile godoc
// @Summary Cria o perfil de responsável por adoção do usuário autenticado
// @Description Cria um perfil de responsável por adoção vinculado ao usuário autenticado. Cada usuário pode ter apenas um perfil. Após criação, o perfil fica aguardando aprovação pelo administrador.
// @Tags Adoção - Perfil do Responsável
// @Accept json
// @Produce json
// @Param input body usecases.CreateAdoptionGuardianProfileInputBody true "Dados do perfil de responsável"
// @Success 201 {object} usecases.CreateAdoptionGuardianProfileOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/guardian-profile [post]
func (h *AdoptionGuardianHandler) CreateAdoptionGuardianProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionGuardianHandler.CreateAdoptionGuardianProfile", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.CreateAdoptionGuardianProfileInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do perfil de responsável",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionGuardianHandler.CreateAdoptionGuardianProfile", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.CreateAdoptionGuardianProfileInput{
		UserID:       userID.(string),
		DisplayName:  inputBody.DisplayName,
		GuardianType: inputBody.GuardianType,
		Document:     inputBody.Document,
		Phone:        inputBody.Phone,
		Whatsapp:     inputBody.Whatsapp,
		About:        inputBody.About,
		CityID:       inputBody.CityID,
		StateID:      inputBody.StateID,
	}

	output, errs := h.AdoptionGuardianFactory.CreateAdoptionGuardianProfile.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// GetMyAdoptionGuardianProfile godoc
// @Summary Retorna o perfil de responsável por adoção do usuário autenticado
// @Description Retorna os dados do perfil de responsável por adoção vinculado ao usuário autenticado, incluindo o status de aprovação.
// @Tags Adoção - Perfil do Responsável
// @Produce json
// @Success 200 {object} usecases.GetMyAdoptionGuardianProfileOutput
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/guardian-profile/me [get]
func (h *AdoptionGuardianHandler) GetMyAdoptionGuardianProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionGuardianHandler.GetMyAdoptionGuardianProfile", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	input := usecases.GetMyAdoptionGuardianProfileInput{
		UserID: userID.(string),
	}

	output, errs := h.AdoptionGuardianFactory.GetMyAdoptionGuardianProfile.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateAdoptionGuardianProfile godoc
// @Summary Atualiza o perfil de responsável por adoção do usuário autenticado
// @Description Atualiza os dados do perfil de responsável por adoção vinculado ao usuário autenticado. Apenas os campos enviados no corpo são atualizados.
// @Tags Adoção - Perfil do Responsável
// @Accept json
// @Produce json
// @Param input body usecases.UpdateAdoptionGuardianProfileInputBody true "Dados a atualizar"
// @Success 200 {object} usecases.UpdateAdoptionGuardianProfileOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/guardian-profile/me [put]
func (h *AdoptionGuardianHandler) UpdateAdoptionGuardianProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionGuardianHandler.UpdateAdoptionGuardianProfile", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.UpdateAdoptionGuardianProfileInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do perfil de responsável",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionGuardianHandler.UpdateAdoptionGuardianProfile", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.UpdateAdoptionGuardianProfileInput{
		UserID:       userID.(string),
		DisplayName:  inputBody.DisplayName,
		GuardianType: inputBody.GuardianType,
		Document:     inputBody.Document,
		Phone:        inputBody.Phone,
		Whatsapp:     inputBody.Whatsapp,
		About:        inputBody.About,
		CityID:       inputBody.CityID,
		StateID:      inputBody.StateID,
	}

	output, errs := h.AdoptionGuardianFactory.UpdateAdoptionGuardianProfile.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
