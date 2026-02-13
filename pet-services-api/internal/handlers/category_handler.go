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

// ListCategories godoc
// @Summary Lista categorias com paginação
// @Tags Categorias
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Param name query string false "Filtro por nome"
// @Success 200 {object} usecases.ListCategoriesOutput
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	ctx := c.Request.Context()

	userType, exists := c.Get("user_type")
	if !exists || userType != "provider" {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Acesso negado",
			Detail: "Acesso permitido apenas para usuários do tipo provider",
		})
		h.Logger.LogError(ctx, "CategoryHandler.ListCategories", problem.Title+": "+problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	page := 1
	pageSize := 10
	name := c.Query("name")
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

	providerID, _ := c.Get("user_id")
	input := usecases.ListCategoriesInput{
		Page:       page,
		PageSize:   pageSize,
		Name:       name,
		ProviderID: providerID.(string),
	}

	output, err := h.CategoryFactory.ListCategories.Execute(ctx, input)
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao listar categorias",
			Detail: err.Error(),
		})
		h.Logger.LogError(ctx, "CategoryHandler.ListCategories", problem.Title+": "+problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
		return
	}

	c.JSON(http.StatusOK, output)
}
