package handlers

import (
	"net/http"

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
