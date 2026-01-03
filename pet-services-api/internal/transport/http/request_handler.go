package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/request"
	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainrequest "github.com/guilherme/pet-services-api/internal/domain/request"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// parseUUIDParamProblems converte parâmetro de rota para UUID, retornando problemas padronizados.
func parseUUIDParamProblems(c *gin.Context, name string, _ string) (uuid.UUID, []exceptions.ProblemDetails) {
	value := c.Param(name)
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "UUID inválido",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}
	}
	return id, nil
}

// RequestHandler expõe endpoints de solicitações de serviço.
type RequestHandler struct {
	uc           factory.RequestUseCases
	providerRepo domainprovider.Repository
	errorService *ErrorService
}

func NewRequestHandler(uc factory.RequestUseCases, providerRepo domainprovider.Repository, errorService *ErrorService) *RequestHandler {
	return &RequestHandler{uc: uc, providerRepo: providerRepo, errorService: errorService}
}

// RegisterRequestRoutes registra rotas autenticadas de solicitações.
func RegisterRequestRoutes(rg *gin.RouterGroup, h *RequestHandler) {
	rg.POST("", h.Create)
	rg.POST(":request_id/accept", h.Accept)
	rg.POST(":request_id/reject", h.Reject)
	rg.POST(":request_id/complete", h.Complete)
	rg.POST(":request_id/cancel", h.Cancel)
	rg.GET(":request_id/status", h.GetStatus)
	rg.GET("/owner", h.ListForOwner)
	rg.GET("/provider", h.ListForProvider)
	rg.GET("/status", h.ListByStatus)
}

type createRequestPayload struct {
	ProviderID    string                `json:"provider_id" validate:"required,uuid"`
	ServiceType   string                `json:"service_type" validate:"required,min=2,max=50"`
	PetName       string                `json:"pet_name" validate:"required,min=1,max=100"`
	PetType       domainrequest.PetType `json:"pet_type" validate:"required"`
	PetBreed      string                `json:"pet_breed" validate:"omitempty,max=100"`
	PetAge        int                   `json:"pet_age" validate:"required,min=0,max=50"`
	PetWeight     float64               `json:"pet_weight" validate:"required,min=0.1,max=100"`
	PetNotes      string                `json:"pet_notes" validate:"omitempty,max=500"`
	PreferredDate string                `json:"preferred_date" validate:"required,datetime=2006-01-02"`
	PreferredTime string                `json:"preferred_time" validate:"required,datetime=15:04"`
	Notes         string                `json:"notes" validate:"omitempty,max=1000"`
}

// Create cria uma solicitação de serviço.
// @Summary Criar solicitação de serviço
// @Tags requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreateRequestInput true "Dados da solicitação"
// @Success 201 {object} map[string]interface{} "ID e status da solicitação criada"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests [post]
func (h *RequestHandler) Create(c *gin.Context) {
	ownerID, err := extractUserID(c)
	if err != nil {
		h.errorService.RespondWithError(c, err, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		h.errorService.RespondWithError(c, errors.New("apenas donos podem criar solicitações"), "forbidden", http.StatusForbidden)
		return
	}

	var req createRequestPayload
	if err := BindAndValidateJSON(c, &req); err != nil {
		h.errorService.RespondWithError(c, err, "invalid_payload", http.StatusBadRequest)
		return
	}

	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		h.errorService.RespondWithError(c, err, "invalid_provider_id", http.StatusBadRequest)
		return
	}

	preferredDate, err := time.Parse(time.RFC3339, req.PreferredDate)
	if err != nil {
		h.errorService.RespondWithError(c, errors.New("use formato ISO8601"), "invalid_preferred_date", http.StatusBadRequest)
		return
	}

	input := request.CreateRequestInput{
		OwnerID:       ownerID,
		ProviderID:    providerID,
		ServiceType:   req.ServiceType,
		PetName:       req.PetName,
		PetType:       req.PetType,
		PetBreed:      req.PetBreed,
		PetAge:        req.PetAge,
		PetWeight:     req.PetWeight,
		PetNotes:      req.PetNotes,
		PreferredDate: preferredDate,
		PreferredTime: req.PreferredTime,
		Notes:         req.Notes,
	}

	created, problems := h.uc.Create.Execute(c.Request.Context(), input)
	if len(problems) > 0 {
		// Busca o status do primeiro problema para resposta
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "create_request_failed", status)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             created.ID,
		"status":         created.Status,
		"preferred_date": created.PreferredDate,
		"preferred_time": created.PreferredTime,
	})
}

// Accept permite ao prestador aceitar uma solicitação.
// @Summary Aceitar solicitação
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} map[string]interface{} "Status de aceitação"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/{request_id}/accept [post]
func (h *RequestHandler) Accept(c *gin.Context) {
	requestID, problems := parseUUIDParamProblems(c, "request_id", "invalid_request_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_request_id", http.StatusBadRequest)
		return
	}
	providerID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider", http.StatusBadRequest)
		return
	}

	problems = h.uc.Accept.Execute(c.Request.Context(), request.AcceptRequestInput{RequestID: requestID, ProviderID: providerID})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "accept_request_failed", status)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "accepted"})
}

// Reject permite ao prestador rejeitar uma solicitação.
// @Summary Rejeitar solicitação
// @Tags requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Param request body object true "Motivo da rejeição"
// @Success 200 {object} map[string]interface{} "Status de rejeição"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/{request_id}/reject [post]
func (h *RequestHandler) Reject(c *gin.Context) {
	requestID, problems := parseUUIDParamProblems(c, "request_id", "invalid_request_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_request_id", http.StatusBadRequest)
		return
	}
	providerID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider", http.StatusBadRequest)
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorService.RespondWithError(c, err, "invalid_payload", http.StatusBadRequest)
		return
	}
	problems = h.uc.Reject.Execute(c.Request.Context(), request.RejectRequestInput{RequestID: requestID, ProviderID: providerID, Reason: req.Reason})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "reject_request_failed", status)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

// Complete marca a solicitação como concluída.
// @Summary Concluir solicitação
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} map[string]interface{} "Status de conclusão"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/{request_id}/complete [post]
func (h *RequestHandler) Complete(c *gin.Context) {
	requestID, problems := parseUUIDParamProblems(c, "request_id", "invalid_request_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_request_id", http.StatusBadRequest)
		return
	}
	providerID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider", http.StatusBadRequest)
		return
	}

	problems = h.uc.Complete.Execute(c.Request.Context(), request.CompleteRequestInput{RequestID: requestID, ProviderID: providerID})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "complete_request_failed", status)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

// Cancel permite ao dono cancelar a solicitação.
// @Summary Cancelar solicitação
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} map[string]interface{} "Status de cancelamento"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/{request_id}/cancel [post]
func (h *RequestHandler) Cancel(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}
	ownerID, err := extractUserID(c)
	if err != nil {
		h.errorService.RespondWithError(c, err, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		h.errorService.RespondWithError(c, errors.New("apenas donos podem cancelar solicitações"), "forbidden", http.StatusForbidden)
		return
	}

	problems := h.uc.Cancel.Execute(c.Request.Context(), request.CancelRequestInput{RequestID: requestID, OwnerID: ownerID})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "cancel_request_failed", status)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

// GetStatus retorna detalhes da solicitação para dono ou prestador vinculado.
// @Summary Consultar status da solicitação
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param request_id path string true "ID da solicitação"
// @Success 200 {object} map[string]interface{} "Detalhes da solicitação"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 403 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/{request_id}/status [get]
func (h *RequestHandler) GetStatus(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}

	userID, _ := extractUserID(c)
	userType := extractUserType(c)
	providerID := uuid.Nil
	if userType == domainuser.UserTypeProvider {
		resolvedID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
		if len(problems) > 0 {
			h.errorService.RespondWithProblems(c, problems, "invalid_provider", http.StatusBadRequest)
			return
		}
		providerID = resolvedID
		userID = uuid.Nil // evita conflitar owner vs provider
	}

	result, problems := h.uc.GetStatus.Execute(c.Request.Context(), request.GetRequestStatusInput{RequestID: requestID, OwnerID: userID, ProviderID: providerID})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "cannot_view_request", status)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               result.ID,
		"status":           result.Status,
		"rejection_reason": result.RejectionReason,
		"provider_id":      result.ProviderID,
		"owner_id":         result.OwnerID,
		"created_at":       result.CreatedAt,
	})
}

// ListForOwner lista solicitações do dono autenticado.
// @Summary Listar solicitações do dono
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} map[string]interface{} "Lista de solicitações"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/owner [get]
func (h *RequestHandler) ListForOwner(c *gin.Context) {
	ownerID, err := extractUserID(c)
	if err != nil {
		h.errorService.RespondWithError(c, err, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		h.errorService.RespondWithError(c, errors.New("apenas donos podem listar suas solicitações"), "forbidden", http.StatusForbidden)
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, problems := h.uc.ListForOwner.Execute(c.Request.Context(), request.ListRequestsForOwnerInput{OwnerID: ownerID, Page: page, Limit: limit})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "list_requests_failed", status)
		return
	}

	resp := mapRequests(items)
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

// ListForProvider lista solicitações do prestador autenticado.
// @Summary Listar solicitações do prestador
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} map[string]interface{} "Lista de solicitações"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/provider [get]
func (h *RequestHandler) ListForProvider(c *gin.Context) {
	providerID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider", http.StatusBadRequest)
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, problems := h.uc.ListForProvider.Execute(c.Request.Context(), request.ListRequestsForProviderInput{ProviderID: providerID, Page: page, Limit: limit})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "list_requests_failed", status)
		return
	}

	resp := mapRequests(items)
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

// ListByStatus lista solicitações por status.
// @Summary Listar solicitações por status
// @Tags requests
// @Security BearerAuth
// @Produce json
// @Param status query string true "Status"
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} map[string]interface{} "Lista de solicitações"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /requests/status [get]
func (h *RequestHandler) ListByStatus(c *gin.Context) {
	status := domainrequest.Status(c.Query("status"))
	if status == "" {
		h.errorService.RespondWithError(c, errors.New("status é obrigatório"), "invalid_status", http.StatusBadRequest)
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, problems := h.uc.ListByStatus.Execute(c.Request.Context(), request.ListRequestsByStatusInput{Status: status, Page: page, Limit: limit})
	if len(problems) > 0 {
		status := http.StatusBadRequest
		for _, p := range problems {
			if p.Status >= 500 {
				status = http.StatusInternalServerError
				break
			} else if p.Status == http.StatusNotFound {
				status = http.StatusNotFound
			}
		}
		h.errorService.RespondWithProblems(c, problems, "list_requests_failed", status)
		return
	}
	resp := mapRequests(items)
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

// mapRequests serializa as solicitações para saída simples.
func mapRequests(requests []*domainrequest.ServiceRequest) []gin.H {
	resp := make([]gin.H, 0, len(requests))
	for _, r := range requests {
		resp = append(resp, gin.H{
			"id":               r.ID,
			"status":           r.Status,
			"service_type":     r.ServiceType,
			"preferred_date":   r.PreferredDate,
			"preferred_time":   r.PreferredTime,
			"provider_id":      r.ProviderID,
			"owner_id":         r.OwnerID,
			"rejection_reason": r.RejectionReason,
			"created_at":       r.CreatedAt,
		})
	}
	return resp
}

// parseUUIDParam converte parâmetro de rota para UUID e retorna se foi válido.
func parseUUIDParam(c *gin.Context, name string, errCode string) (uuid.UUID, bool) {
	value := c.Param(name)
	id, err := uuid.Parse(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errCode, err.Error()))
		return uuid.Nil, false
	}
	return id, true
}
