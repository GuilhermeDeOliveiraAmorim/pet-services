package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/request"
	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainrequest "github.com/guilherme/pet-services-api/internal/domain/request"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// RequestHandler expõe endpoints de solicitações de serviço.
type RequestHandler struct {
	uc           factory.RequestUseCases
	providerRepo domainprovider.Repository
}

func NewRequestHandler(uc factory.RequestUseCases, providerRepo domainprovider.Repository) *RequestHandler {
	return &RequestHandler{uc: uc, providerRepo: providerRepo}
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
	ProviderID    string                `json:"provider_id"`
	ServiceType   string                `json:"service_type"`
	PetName       string                `json:"pet_name"`
	PetType       domainrequest.PetType `json:"pet_type"`
	PetBreed      string                `json:"pet_breed"`
	PetAge        int                   `json:"pet_age"`
	PetWeight     float64               `json:"pet_weight"`
	PetNotes      string                `json:"pet_notes"`
	PreferredDate string                `json:"preferred_date"` // ISO date
	PreferredTime string                `json:"preferred_time"` // HH:MM
	Notes         string                `json:"notes"`
}

func (h *RequestHandler) Create(c *gin.Context) {
	ownerID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		c.JSON(http.StatusForbidden, errorResponse("forbidden", "apenas donos podem criar solicitações"))
		return
	}

	var req createRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}

	preferredDate, err := time.Parse(time.RFC3339, req.PreferredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_preferred_date", "use formato ISO8601"))
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

	created, err := h.uc.Create.Execute(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		case errors.Is(err, domainprovider.ErrProviderNotActive), errors.Is(err, domainrequest.ErrInvalidPreferredDate), errors.Is(err, domainrequest.ErrInvalidServiceType):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_request", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("create_request_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             created.ID,
		"status":         created.Status,
		"preferred_date": created.PreferredDate,
		"preferred_time": created.PreferredTime,
	})
}

func (h *RequestHandler) Accept(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}
	providerID, ok := providerIDFromContext(c, h.providerRepo, true)
	if !ok {
		return
	}

	if err := h.uc.Accept.Execute(c.Request.Context(), request.AcceptRequestInput{RequestID: requestID, ProviderID: providerID}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainrequest.ErrRequestNotFound):
			c.JSON(http.StatusNotFound, errorResponse("request_or_provider_not_found", err.Error()))
		case errors.Is(err, domainprovider.ErrProviderNotActive), errors.Is(err, domainrequest.ErrInvalidStatusTransition):
			c.JSON(http.StatusBadRequest, errorResponse("cannot_accept", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("accept_request_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "accepted"})
}

func (h *RequestHandler) Reject(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}
	providerID, ok := providerIDFromContext(c, h.providerRepo, true)
	if !ok {
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}
	if err := h.uc.Reject.Execute(c.Request.Context(), request.RejectRequestInput{RequestID: requestID, ProviderID: providerID, Reason: req.Reason}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainrequest.ErrRequestNotFound):
			c.JSON(http.StatusNotFound, errorResponse("request_or_provider_not_found", err.Error()))
		case errors.Is(err, domainprovider.ErrProviderNotActive), errors.Is(err, domainrequest.ErrInvalidStatusTransition):
			c.JSON(http.StatusBadRequest, errorResponse("cannot_reject", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("reject_request_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

func (h *RequestHandler) Complete(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}
	providerID, ok := providerIDFromContext(c, h.providerRepo, true)
	if !ok {
		return
	}

	if err := h.uc.Complete.Execute(c.Request.Context(), request.CompleteRequestInput{RequestID: requestID, ProviderID: providerID}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainrequest.ErrRequestNotFound):
			c.JSON(http.StatusNotFound, errorResponse("request_or_provider_not_found", err.Error()))
		case errors.Is(err, domainrequest.ErrInvalidStatusTransition):
			c.JSON(http.StatusBadRequest, errorResponse("cannot_complete", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("complete_request_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

func (h *RequestHandler) Cancel(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}
	ownerID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		c.JSON(http.StatusForbidden, errorResponse("forbidden", "apenas donos podem cancelar solicitações"))
		return
	}

	if err := h.uc.Cancel.Execute(c.Request.Context(), request.CancelRequestInput{RequestID: requestID, OwnerID: ownerID}); err != nil {
		switch {
		case errors.Is(err, domainrequest.ErrRequestNotFound):
			c.JSON(http.StatusNotFound, errorResponse("request_not_found", err.Error()))
		case errors.Is(err, domainrequest.ErrInvalidStatusTransition):
			c.JSON(http.StatusBadRequest, errorResponse("cannot_cancel", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("cancel_request_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *RequestHandler) GetStatus(c *gin.Context) {
	requestID, ok := parseUUIDParam(c, "request_id", "invalid_request_id")
	if !ok {
		return
	}

	userID, _ := extractUserID(c)
	userType := extractUserType(c)
	providerID := uuid.Nil
	if userType == domainuser.UserTypeProvider {
		resolvedID, ok := providerIDFromContext(c, h.providerRepo, true)
		if !ok {
			return
		}
		providerID = resolvedID
		userID = uuid.Nil // evita conflitar owner vs provider
	}

	result, err := h.uc.GetStatus.Execute(c.Request.Context(), request.GetRequestStatusInput{RequestID: requestID, OwnerID: userID, ProviderID: providerID})
	if err != nil {
		switch {
		case errors.Is(err, domainrequest.ErrRequestNotFound):
			c.JSON(http.StatusNotFound, errorResponse("request_not_found", err.Error()))
		default:
			c.JSON(http.StatusForbidden, errorResponse("cannot_view_request", err.Error()))
		}
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

func (h *RequestHandler) ListForOwner(c *gin.Context) {
	ownerID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		c.JSON(http.StatusForbidden, errorResponse("forbidden", "apenas donos podem listar suas solicitações"))
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, err := h.uc.ListForOwner.Execute(c.Request.Context(), request.ListRequestsForOwnerInput{OwnerID: ownerID, Page: page, Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("list_requests_failed", err.Error()))
		return
	}

	resp := mapRequests(items)
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

func (h *RequestHandler) ListForProvider(c *gin.Context) {
	providerID, ok := providerIDFromContext(c, h.providerRepo, true)
	if !ok {
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, err := h.uc.ListForProvider.Execute(c.Request.Context(), request.ListRequestsForProviderInput{ProviderID: providerID, Page: page, Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("list_requests_failed", err.Error()))
		return
	}

	resp := mapRequests(items)
	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

func (h *RequestHandler) ListByStatus(c *gin.Context) {
	status := domainrequest.Status(c.Query("status"))
	if status == "" {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_status", "status é obrigatório"))
		return
	}
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	items, total, err := h.uc.ListByStatus.Execute(c.Request.Context(), request.ListRequestsByStatusInput{Status: status, Page: page, Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("list_requests_failed", err.Error()))
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
