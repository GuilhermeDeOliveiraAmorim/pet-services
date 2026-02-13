package handlers

import (
	"net/http"

	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceFactory *factories.ServiceFactory
	Logger         logging.LoggerInterface
}

func NewServiceHandler(factory *factories.ServiceFactory, logger logging.LoggerInterface) *ServiceHandler {
	return &ServiceHandler{
		ServiceFactory: factory,
		Logger:         logger,
	}
}

// AddService godoc
// @Summary Adiciona um novo serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param input body usecases.AddServiceInputBody true "Dados do serviço"
// @Success 201 {object} usecases.AddServiceOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /services [post]
func (h *ServiceHandler) AddService(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddService", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddServiceInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do serviço",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddService", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddServiceInput{
		UserID:       userID.(string),
		Name:         inputBody.Name,
		Description:  inputBody.Description,
		Price:        inputBody.Price,
		PriceMinimum: inputBody.PriceMinimum,
		PriceMaximum: inputBody.PriceMaximum,
		Duration:     inputBody.Duration,
	}

	output, errs := h.ServiceFactory.AddService.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
