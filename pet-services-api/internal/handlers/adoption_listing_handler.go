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

type AdoptionListingHandler struct {
	AdoptionListingFactory *factories.AdoptionListingFactory
	Logger                 logging.LoggerInterface
}

func NewAdoptionListingHandler(factory *factories.AdoptionListingFactory, logger logging.LoggerInterface) *AdoptionListingHandler {
	return &AdoptionListingHandler{
		AdoptionListingFactory: factory,
		Logger:                 logger,
	}
}

// CreateAdoptionListing godoc
// @Summary Cria um anúncio de adoção
// @Description Cria um novo anúncio de adoção vinculado ao perfil de responsável aprovado do usuário autenticado. O anúncio é criado com status 'draft'.
// @Tags Adoção - Anúncios
// @Accept json
// @Produce json
// @Param input body usecases.CreateAdoptionListingInputBody true "Dados do anúncio"
// @Success 201 {object} usecases.CreateAdoptionListingOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings [post]
func (h *AdoptionListingHandler) CreateAdoptionListing(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.CreateAdoptionListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	guardianProfileID, exists := c.Get("guardian_profile_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Perfil não aprovado",
			Detail: "Apenas responsáveis com perfil aprovado podem criar anúncios de adoção",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.CreateAdoptionListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	var inputBody usecases.CreateAdoptionListingInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do anúncio",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.CreateAdoptionListing", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.CreateAdoptionListingInput{
		UserID:                         userID.(string),
		GuardianProfileID:              guardianProfileID.(string),
		CreateAdoptionListingInputBody: inputBody,
	}

	output, errs := h.AdoptionListingFactory.CreateAdoptionListing.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// UpdateAdoptionListing godoc
// @Summary Atualiza um anúncio de adoção
// @Description Atualiza os dados de um anúncio de adoção em status 'draft' ou 'paused'. Apenas o responsável pelo anúncio pode editá-lo.
// @Tags Adoção - Anúncios
// @Accept json
// @Produce json
// @Param listing_id path string true "ID do anúncio"
// @Param input body usecases.UpdateAdoptionListingInputBody true "Campos a atualizar"
// @Success 200 {object} usecases.UpdateAdoptionListingOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings/{listing_id} [put]
func (h *AdoptionListingHandler) UpdateAdoptionListing(c *gin.Context) {
	ctx := c.Request.Context()

	guardianProfileID, exists := c.Get("guardian_profile_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Perfil não aprovado",
			Detail: "Apenas responsáveis com perfil aprovado podem editar anúncios de adoção",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.UpdateAdoptionListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	listingID := c.Param("listing_id")
	if listingID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do anúncio ausente",
			Detail: "O ID do anúncio é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.UpdateAdoptionListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.UpdateAdoptionListingInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do anúncio",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.UpdateAdoptionListing", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.UpdateAdoptionListingInput{
		ListingID:                      listingID,
		GuardianProfileID:              guardianProfileID.(string),
		UpdateAdoptionListingInputBody: inputBody,
	}

	output, errs := h.AdoptionListingFactory.UpdateAdoptionListing.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ChangeAdoptionListingStatus godoc
// @Summary Altera o status de um anúncio de adoção
// @Description Altera o status de um anúncio. Ações disponíveis: publish (draft/paused → published), pause (published → paused), archive (qualquer → archived). Apenas o responsável pode realizar esta ação.
// @Tags Adoção - Anúncios
// @Produce json
// @Param listing_id path string true "ID do anúncio"
// @Param action path string true "Ação" Enums(publish, pause, archive)
// @Success 200 {object} usecases.ChangeAdoptionListingStatusOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings/{listing_id}/{action} [patch]
func (h *AdoptionListingHandler) ChangeAdoptionListingStatus(c *gin.Context) {
	ctx := c.Request.Context()

	guardianProfileID, exists := c.Get("guardian_profile_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Perfil não aprovado",
			Detail: "Apenas responsáveis com perfil aprovado podem gerenciar anúncios de adoção",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.ChangeAdoptionListingStatus", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	listingID := c.Param("listing_id")
	action := c.Param("action")

	if listingID == "" || action == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Parâmetros ausentes",
			Detail: "O ID do anúncio e a ação são obrigatórios",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.ChangeAdoptionListingStatus", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.ChangeAdoptionListingStatusInput{
		ListingID:         listingID,
		GuardianProfileID: guardianProfileID.(string),
		Action:            action,
	}

	output, errs := h.AdoptionListingFactory.ChangeAdoptionListingStatus.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListPublicAdoptionListings godoc
// @Summary Lista anúncios públicos de adoção
// @Description Retorna a lista paginada de anúncios de adoção com status 'published'. Suporta filtros por sexo, porte, faixa etária, cidade e estado.
// @Tags Adoção - Público
// @Produce json
// @Param page query int false "Página (padrão: 1)"
// @Param page_size query int false "Itens por página (padrão: 12, máx: 50)"
// @Param sex query string false "Sexo do animal" Enums(male, female)
// @Param size query string false "Porte do animal" Enums(small, medium, large)
// @Param age_group query string false "Faixa etária" Enums(puppy, adult, senior)
// @Param city_id query string false "ID da cidade"
// @Param state_id query string false "ID do estado"
// @Success 200 {object} usecases.ListPublicAdoptionListingsOutput
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /adoption/listings [get]
func (h *AdoptionListingHandler) ListPublicAdoptionListings(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))

	input := usecases.ListPublicAdoptionListingsInput{
		Page:     page,
		PageSize: pageSize,
		Sex:      c.Query("sex"),
		Size:     c.Query("size"),
		AgeGroup: c.Query("age_group"),
		CityID:   c.Query("city_id"),
		StateID:  c.Query("state_id"),
	}

	output, errs := h.AdoptionListingFactory.ListPublicAdoptionListings.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetPublicAdoptionListing godoc
// @Summary Retorna o detalhe de um anúncio público de adoção
// @Description Retorna os dados completos de um anúncio de adoção com status 'published', incluindo dados do pet e do responsável.
// @Tags Adoção - Público
// @Produce json
// @Param listing_id path string true "ID do anúncio"
// @Success 200 {object} usecases.GetPublicAdoptionListingOutput
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /adoption/listings/{listing_id} [get]
func (h *AdoptionListingHandler) GetPublicAdoptionListing(c *gin.Context) {
	ctx := c.Request.Context()

	listingID := c.Param("listing_id")
	if listingID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do anúncio ausente",
			Detail: "O ID do anúncio é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.GetPublicAdoptionListing", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	output, errs := h.AdoptionListingFactory.GetPublicAdoptionListing.Execute(ctx, usecases.GetPublicAdoptionListingInput{
		ListingID: listingID,
	})
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListMyAdoptionListings godoc
// @Summary Lista os anúncios de adoção do responsável autenticado
// @Description Retorna a lista paginada de anúncios de adoção vinculados ao perfil de responsável do usuário autenticado, em todos os status.
// @Tags Adoção - Anúncios
// @Produce json
// @Param page query int false "Página (padrão: 1)"
// @Param page_size query int false "Itens por página (padrão: 10, máx: 50)"
// @Success 200 {object} usecases.ListMyAdoptionListingsOutput
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings/me [get]
func (h *AdoptionListingHandler) ListMyAdoptionListings(c *gin.Context) {
	ctx := c.Request.Context()

	guardianProfileID, exists := c.Get("guardian_profile_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Perfil não aprovado",
			Detail: "Apenas responsáveis com perfil aprovado podem listar seus anúncios",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.ListMyAdoptionListings", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	output, errs := h.AdoptionListingFactory.ListMyAdoptionListings.Execute(ctx, usecases.ListMyAdoptionListingsInput{
		GuardianProfileID: guardianProfileID.(string),
		Page:              page,
		PageSize:          pageSize,
	})
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// MarkAdoptionListingAsAdopted godoc
// @Summary Marca um anúncio como adotado
// @Description Marca um anúncio como adotado. Apenas o responsável proprietário do anúncio pode executar esta ação.
// @Tags Adoção - Anúncios
// @Accept json
// @Produce json
// @Param id path string true "ID do anúncio"
// @Success 200 {object} usecases.MarkAdoptionListingAsAdoptedOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /adoption/listings/{id}/mark-adopted [post]
func (h *AdoptionListingHandler) MarkAdoptionListingAsAdopted(c *gin.Context) {
	ctx := c.Request.Context()

	listingID := c.Param("id")
	guardianProfileID, exists := c.Get("guardian_profile_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Perfil não aprovado",
			Detail: "Apenas responsáveis com perfil aprovado podem marcar um anúncio como adotado",
		})
		h.Logger.LogBadRequest(ctx, "AdoptionListingHandler.MarkAdoptionListingAsAdopted", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	output, errs := h.AdoptionListingFactory.MarkAdoptionListingAsAdopted.Execute(ctx, usecases.MarkAdoptionListingAsAdoptedInput{
		ListingID:         listingID,
		GuardianProfileID: guardianProfileID.(string),
	})
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
