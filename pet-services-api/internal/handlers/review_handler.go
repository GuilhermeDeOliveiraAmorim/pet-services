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

type ReviewHandler struct {
	ReviewFactory *factories.ReviewFactory
	Logger        logging.LoggerInterface
}

func NewReviewHandler(factory *factories.ReviewFactory, logger logging.LoggerInterface) *ReviewHandler {
	return &ReviewHandler{
		ReviewFactory: factory,
		Logger:        logger,
	}
}

// CreateReview godoc
// @Summary Cria um review para provedor
// @Tags Reviews
// @Accept json
// @Produce json
// @Param provider_id path string true "ID do provedor"
// @Param input body usecases.CreateReviewInputBody true "Dados do review"
// @Success 201 {object} usecases.CreateReviewOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /providers/{provider_id}/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ReviewHandler.CreateReview", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	providerID := c.Param("provider_id")
	if providerID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do provedor ausente",
			Detail: "O ID do provedor é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "ReviewHandler.CreateReview", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.CreateReviewInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do review",
		})
		h.Logger.LogBadRequest(ctx, "ReviewHandler.CreateReview", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.CreateReviewInput{
		UserID:     userID.(string),
		ProviderID: providerID,
		Rating:     inputBody.Rating,
		Comment:    inputBody.Comment,
	}

	output, errs := h.ReviewFactory.CreateReview.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// ListReviews godoc
// @Summary Lista reviews por provedor ou usuário
// @Tags Reviews
// @Accept json
// @Produce json
// @Param provider_id query string false "Filtro por provedor"
// @Param user_id query string false "Filtro por usuário"
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Success 200 {object} usecases.ListReviewsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Router /reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	ctx := c.Request.Context()

	providerID := c.Query("provider_id")
	userID := c.Query("user_id")

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 {
			pageSize = val
		}
	}

	input := usecases.ListReviewsInput{
		ProviderID: providerID,
		UserID:     userID,
		Page:       page,
		PageSize:   pageSize,
	}

	output, errs := h.ReviewFactory.ListReviews.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
