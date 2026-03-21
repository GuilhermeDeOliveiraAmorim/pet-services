package handlers

import (
	"net/http"
	"strconv"

	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type AdoptionApplicationHandler struct {
	AdoptionApplicationFactory *factories.AdoptionApplicationFactory
	Logger                     logging.LoggerInterface
}

func NewAdoptionApplicationHandler(factory *factories.AdoptionApplicationFactory, logger logging.LoggerInterface) *AdoptionApplicationHandler {
	return &AdoptionApplicationHandler{
		AdoptionApplicationFactory: factory,
		Logger:                     logger,
	}
}

// CreateAdoptionApplication godoc
// @Summary Cria uma candidatura de adoção
// @Description Submete uma candidatura para adotar um pet a partir de um anúncio publicado. O usuário autenticado é registrado como candidato.
// @Tags Adoção - Candidaturas
// @Accept json
// @Produce json
// @Param input body usecases.CreateAdoptionApplicationInputBody true "Dados da candidatura"
// @Success 201 {object} usecases.CreateAdoptionApplicationOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/applications [post]
func (h *AdoptionApplicationHandler) CreateAdoptionApplication(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.CreateAdoptionApplication", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.CreateAdoptionApplicationInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da candidatura",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.CreateAdoptionApplication", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.CreateAdoptionApplicationInput{
		UserID:                             userID.(string),
		CreateAdoptionApplicationInputBody: inputBody,
	}

	output, errs := h.AdoptionApplicationFactory.CreateAdoptionApplication.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// ListMyAdoptionApplications godoc
// @Summary Lista as minhas candidaturas de adoção
// @Description Retorna a lista paginada de candidaturas do usuário autenticado em todos os estágios.
// @Tags Adoção - Candidaturas
// @Produce json
// @Param page query int false "Página (padrão: 1)"
// @Param page_size query int false "Itens por página (padrão: 10, máx: 50)"
// @Success 200 {object} usecases.ListMyAdoptionApplicationsOutput
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/applications/me [get]
func (h *AdoptionApplicationHandler) ListMyAdoptionApplications(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.ListMyAdoptionApplications", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	input := usecases.ListMyAdoptionApplicationsInput{
		UserID:   userID.(string),
		Page:     page,
		PageSize: pageSize,
	}

	output, errs := h.AdoptionApplicationFactory.ListMyAdoptionApplications.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListAdoptionApplicationsByListing godoc
// @Summary Lista as candidaturas de um anúncio
// @Description Retorna a lista paginada de candidaturas para um anúncio específico. Apenas o responsável pelo anúncio ou admin podem acessar.
// @Tags Adoção - Candidaturas
// @Produce json
// @Param listing_id path string true "ID do anúncio"
// @Param page query int false "Página (padrão: 1)"
// @Param page_size query int false "Itens por página (padrão: 10, máx: 50)"
// @Success 200 {object} usecases.ListAdoptionApplicationsByListingOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings/{listing_id}/applications [get]
func (h *AdoptionApplicationHandler) ListAdoptionApplicationsByListing(c *gin.Context) {
	ctx := c.Request.Context()

	listingID := c.Param("listing_id")
	if listingID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do anúncio ausente",
			Detail: "O ID do anúncio é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.ListAdoptionApplicationsByListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	input := usecases.ListAdoptionApplicationsByListingInput{
		ListingID: listingID,
		Page:      page,
		PageSize:  pageSize,
	}

	output, errs := h.AdoptionApplicationFactory.ListAdoptionApplicationsByListing.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ReviewAdoptionApplication godoc
// @Summary Analisa e muda o status de uma candidatura
// @Description Permite o responsável pelo anúncio ou admin mover a candidatura para analysis, entrevista, aprovação ou rejeição. A ação determina a transição de status.
// @Tags Adoção - Candidaturas
// @Accept json
// @Produce json
// @Param application_id path string true "ID da candidatura"
// @Param input body usecases.ReviewAdoptionApplicationInputBody true "Ação e notas"
// @Success 200 {object} usecases.ReviewAdoptionApplicationOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/applications/{application_id}/review [post]
func (h *AdoptionApplicationHandler) ReviewAdoptionApplication(c *gin.Context) {
	ctx := c.Request.Context()

	reviewerID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.ReviewAdoptionApplication", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	applicationID := c.Param("application_id")
	if applicationID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da candidatura ausente",
			Detail: "O ID da candidatura é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.ReviewAdoptionApplication", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.ReviewAdoptionApplicationInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da análise",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.ReviewAdoptionApplication", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.ReviewAdoptionApplicationInput{
		ApplicationID:                      applicationID,
		ReviewerID:                         reviewerID.(string),
		ReviewAdoptionApplicationInputBody: inputBody,
	}

	output, errs := h.AdoptionApplicationFactory.ReviewAdoptionApplication.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// WithdrawAdoptionApplication godoc
// @Summary Retira uma candidatura de adoção
// @Description Permite que o candidato retire sua candidatura, desde que ela ainda não tenha sido aprovada, rejeitada ou retirada.
// @Tags Adoção - Candidaturas
// @Produce json
// @Param application_id path string true "ID da candidatura"
// @Success 200 {object} usecases.WithdrawAdoptionApplicationOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/applications/{application_id}/withdraw [post]
func (h *AdoptionApplicationHandler) WithdrawAdoptionApplication(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.WithdrawAdoptionApplication", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	applicationID := c.Param("application_id")
	if applicationID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da candidatura ausente",
			Detail: "O ID da candidatura é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionApplicationHandler.WithdrawAdoptionApplication", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.WithdrawAdoptionApplicationInput{
		ApplicationID: applicationID,
		UserID:        userID.(string),
	}

	output, errs := h.AdoptionApplicationFactory.WithdrawAdoptionApplication.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
