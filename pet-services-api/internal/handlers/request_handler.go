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

type RequestHandler struct {
	RequestFactory *factories.RequestFactory
	Logger         logging.LoggerInterface
}

func NewRequestHandler(factory *factories.RequestFactory, logger logging.LoggerInterface) *RequestHandler {
	return &RequestHandler{
		RequestFactory: factory,
		Logger:         logger,
	}
}

// AddRequest godoc
// @Summary Cria uma solicitação de serviço
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param input body usecases.AddRequestInputBody true "Dados da solicitação"
// @Success 201 {object} usecases.AddRequestOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests [post]
func (h *RequestHandler) AddRequest(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.AddRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddRequestInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da solicitação",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.AddRequest", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddRequestInput{
		UserID:    userID.(string),
		ServiceID: inputBody.ServiceID,
		PetID:     inputBody.PetID,
		Notes:     inputBody.Notes,
	}

	output, errs := h.RequestFactory.AddRequest.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// ListRequests godoc
// @Summary Lista solicitações de serviço
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param status query string false "Filtro por status (pending, accepted, rejected, completed)"
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Success 200 {object} usecases.ListRequestsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests [get]
func (h *RequestHandler) ListRequests(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.ListRequests", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	page := 1
	pageSize := 10
	status := c.Query("status")

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

	input := usecases.ListRequestsInput{
		UserID:   userID.(string),
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	output, errs := h.RequestFactory.ListRequests.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetRequest godoc
// @Summary Obtém detalhes de uma solicitação
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} usecases.GetRequestOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests/{request_id} [get]
func (h *RequestHandler) GetRequest(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.GetRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da solicitação ausente",
			Detail: "O ID da solicitação é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.GetRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.GetRequestInput{
		UserID:    userID.(string),
		RequestID: requestID,
	}

	output, errs := h.RequestFactory.GetRequest.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// AcceptRequest godoc
// @Summary Aceita uma solicitação de serviço
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} usecases.AcceptRequestOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests/{request_id}/accept [patch]
func (h *RequestHandler) AcceptRequest(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.AcceptRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da solicitação ausente",
			Detail: "O ID da solicitação é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.AcceptRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AcceptRequestInput{
		UserID:    userID.(string),
		RequestID: requestID,
	}

	output, errs := h.RequestFactory.AcceptRequest.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// CompleteRequest godoc
// @Summary Conclui uma solicitação de serviço
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} usecases.CompleteRequestOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests/{request_id}/complete [patch]
func (h *RequestHandler) CompleteRequest(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.CompleteRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da solicitação ausente",
			Detail: "O ID da solicitação é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.CompleteRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.CompleteRequestInput{
		UserID:    userID.(string),
		RequestID: requestID,
	}

	output, errs := h.RequestFactory.CompleteRequest.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// RejectRequest godoc
// @Summary Rejeita uma solicitação de serviço
// @Tags Solicitações
// @Accept json
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Param input body usecases.RejectRequestInputBody true "Motivo da rejeição"
// @Success 200 {object} usecases.RejectRequestOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /requests/{request_id}/reject [patch]
func (h *RequestHandler) RejectRequest(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.RejectRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da solicitação ausente",
			Detail: "O ID da solicitação é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.RejectRequest", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.RejectRequestInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do motivo da rejeição",
		})
		h.Logger.LogBadRequest(ctx, "RequestHandler.RejectRequest", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.RejectRequestInput{
		UserID:    userID.(string),
		RequestID: requestID,
		Reason:    inputBody.Reason,
	}

	output, errs := h.RequestFactory.RejectRequest.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
