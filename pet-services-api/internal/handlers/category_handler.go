package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryFactory *factories.CategoryFactory
	Logger          logging.LoggerInterface
}

func NewCategoryHandler(factory *factories.CategoryFactory, logger logging.LoggerInterface) *CategoryHandler {
	return &CategoryHandler{
		CategoryFactory: factory,
		Logger:          logger,
	}
}

// CreateCategory godoc
// @Summary Cria uma nova categoria de serviço
// @Tags Categorias
// @Accept json
// @Produce json
// @Param input body usecases.CreateCategoryInput true "Dados da categoria"
// @Success 201 {object} usecases.CreateCategoryOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	userType, exists := c.Get("user_type")
	if !exists || userType != "admin" {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Acesso negado",
			Detail: "Apenas administradores podem criar categorias",
		})
		h.Logger.LogError(ctx, "CategoryHandler.CreateCategory", problem.Title+": "+problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	var input usecases.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da categoria",
		})
		h.Logger.LogBadRequest(ctx, "CategoryHandler.CreateCategory", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input.IsAdmin = true // Garantido pelo middleware

	output, errs := h.CategoryFactory.CreateCategory.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
