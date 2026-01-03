package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/review"
	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainreview "github.com/guilherme/pet-services-api/internal/domain/review"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

func (s *ErrorService) ListReviewsErrorToDTO(err error) (int, ErrorDTO) {
	switch {
	case errors.Is(err, domainreview.ErrReviewNotFound):
		return http.StatusNotFound, ErrorDTO{Code: "reviews_not_found", Message: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorDTO{Code: "list_reviews_failed", Message: err.Error()}
	}
}

// ErrorDTO padroniza resposta de erro.
// ErrorDTO já está definido em error_service.go

func (s *ErrorService) ReviewErrorToDTO(err error) (int, ErrorDTO) {
	switch {
	case errors.Is(err, domainreview.ErrInvalidRating):
		return http.StatusBadRequest, ErrorDTO{Code: "invalid_rating", Message: err.Error()}
	case errors.Is(err, domainreview.ErrReviewAlreadyExists):
		return http.StatusConflict, ErrorDTO{Code: "review_exists", Message: err.Error()}
	case errors.Is(err, domainreview.ErrRequestNotCompleted):
		return http.StatusBadRequest, ErrorDTO{Code: "request_not_completed", Message: err.Error()}
	default:
		return http.StatusBadRequest, ErrorDTO{Code: "submit_review_failed", Message: err.Error()}
	}
}

// ...existing code...

// ReviewHandler expõe endpoints de avaliações.
type ReviewHandler struct {
	uc           factory.ReviewUseCases
	providerRepo domainprovider.Repository
	errorService *ErrorService
}

func NewReviewHandler(uc factory.ReviewUseCases, providerRepo domainprovider.Repository, errorService *ErrorService) *ReviewHandler {
	return &ReviewHandler{uc: uc, providerRepo: providerRepo, errorService: errorService}
}

// RegisterReviewRoutes registra rotas autenticadas para avaliações.
func RegisterReviewRoutes(rg *gin.RouterGroup, h *ReviewHandler) {
	rg.POST("", h.Submit)
	rg.GET("/provider/:provider_id", h.ListForProvider)
}

// submitReviewRequest estrutura de requisição para submeter avaliação.
type submitReviewRequest struct {
	RequestID string `json:"request_id" validate:"required,uuid"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" validate:"omitempty,max=1000"`
}

// Submit cria avaliação para solicitação concluída.
// @Summary Submit review
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Review payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reviews [post]
func (h *ReviewHandler) Submit(c *gin.Context) {
	ownerID, err := extractUserID(c)
	if err != nil {
		h.errorService.RespondWithError(c, err, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		h.errorService.RespondWithError(c, errors.New("apenas donos podem registrar avaliações"), "forbidden", http.StatusForbidden)
		return
	}
	var req submitReviewRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		h.errorService.RespondWithError(c, err, "invalid_payload", http.StatusBadRequest)
		return
	}
	reqID, err := uuid.Parse(req.RequestID)
	if err != nil {
		h.errorService.RespondWithError(c, err, "invalid_request_id", http.StatusBadRequest)
		return
	}

	reviewCreated, err := h.uc.Submit.Execute(c.Request.Context(), review.SubmitReviewInput{
		RequestID: reqID,
		OwnerID:   ownerID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		status, dto := h.errorService.ReviewErrorToDTO(err)
		h.errorService.RespondWithError(c, errors.New(dto.Message), dto.Code, status)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          reviewCreated.ID,
		"provider_id": reviewCreated.ProviderID,
		"rating":      reviewCreated.Rating,
		"comment":     reviewCreated.Comment,
	})
}

// ListForProvider lista avaliações de um prestador.
// @Summary List provider reviews
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /reviews/provider/{provider_id} [get]
func (h *ReviewHandler) ListForProvider(c *gin.Context) {
	providerID, ok := parseUUIDParam(c, "provider_id", "invalid_provider_id")
	if !ok {
		return
	}

	// Se o usuário autenticado for prestador, ele só pode listar o próprio perfil.
	if extractUserType(c) == domainuser.UserTypeProvider {
		authProviderID, ok := providerIDFromContext(c, h.providerRepo, true)
		if !ok {
			return
		}
		if authProviderID != providerID {
			h.errorService.RespondWithError(c, errors.New("não autorizado a listar avaliações de outro prestador"), "forbidden", http.StatusForbidden)
			return
		}
	}

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	reviews, total, err := h.uc.ListForProvider.Execute(c.Request.Context(), review.ListReviewsForProviderInput{ProviderID: providerID, Page: page, Limit: limit})
	if err != nil {
		status, dto := h.errorService.ListReviewsErrorToDTO(err)
		h.errorService.RespondWithError(c, errors.New(dto.Message), dto.Code, status)
		return
	}

	resp := make([]gin.H, 0, len(reviews))
	for _, r := range reviews {
		resp = append(resp, gin.H{
			"id":         r.ID,
			"rating":     r.Rating,
			"comment":    r.Comment,
			"owner_id":   r.OwnerID,
			"request_id": r.RequestID,
			"created_at": r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}
