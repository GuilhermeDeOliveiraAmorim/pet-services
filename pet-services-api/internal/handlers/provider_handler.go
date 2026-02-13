package handlers

import (
	"net/http"

	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	ProviderFactory *factories.ProviderFactory
	Logger          logging.LoggerInterface
}

func NewProviderHandler(factory *factories.ProviderFactory, logger logging.LoggerInterface) *ProviderHandler {
	return &ProviderHandler{
		ProviderFactory: factory,
		Logger:          logger,
	}
}

// AddProvider godoc
// @Summary Cria um provedor para o usuário autenticado
// @Tags Provedores
// @Accept json
// @Produce json
// @Param input body usecases.AddProviderInputBody true "Dados do provedor"
// @Success 201 {object} usecases.AddProviderOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /providers [post]
func (h *ProviderHandler) AddProvider(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ProviderHandler.AddProvider", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddProviderInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do provedor",
		})
		h.Logger.LogBadRequest(ctx, "ProviderHandler.AddProvider", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddProviderInput{
		UserID:       userID.(string),
		BusinessName: inputBody.BusinessName,
		Description:  inputBody.Description,
		PriceRange:   inputBody.PriceRange,
		Address:      inputBody.Address,
	}

	output, errs := h.ProviderFactory.AddProvider.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
