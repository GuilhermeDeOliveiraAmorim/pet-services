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

// ReviewHandler expõe endpoints de avaliações.
type ReviewHandler struct {
	uc           factory.ReviewUseCases
	providerRepo domainprovider.Repository
}

func NewReviewHandler(uc factory.ReviewUseCases, providerRepo domainprovider.Repository) *ReviewHandler {
	return &ReviewHandler{uc: uc, providerRepo: providerRepo}
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
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		c.JSON(http.StatusForbidden, errorResponse("forbidden", "apenas donos podem registrar avaliações"))
		return
	}
	var req submitReviewRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}
	reqID, err := uuid.Parse(req.RequestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_request_id", err.Error()))
		return
	}

	reviewCreated, err := h.uc.Submit.Execute(c.Request.Context(), review.SubmitReviewInput{
		RequestID: reqID,
		OwnerID:   ownerID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainreview.ErrInvalidRating):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_rating", err.Error()))
		case errors.Is(err, domainreview.ErrReviewAlreadyExists):
			c.JSON(http.StatusConflict, errorResponse("review_exists", err.Error()))
		case errors.Is(err, domainreview.ErrRequestNotCompleted):
			c.JSON(http.StatusBadRequest, errorResponse("request_not_completed", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("submit_review_failed", err.Error()))
		}
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
			c.JSON(http.StatusForbidden, errorResponse("forbidden", "não autorizado a listar avaliações de outro prestador"))
			return
		}
	}

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	reviews, total, err := h.uc.ListForProvider.Execute(c.Request.Context(), review.ListReviewsForProviderInput{ProviderID: providerID, Page: page, Limit: limit})
	if err != nil {
		switch {
		case errors.Is(err, domainreview.ErrReviewNotFound):
			c.JSON(http.StatusNotFound, errorResponse("reviews_not_found", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("list_reviews_failed", err.Error()))
		}
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
