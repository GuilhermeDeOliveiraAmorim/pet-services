package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/review"
	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainreview "github.com/guilherme/pet-services-api/internal/domain/review"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// ListReviewsResponseDTO representa a resposta da listagem de avaliações de um prestador.
type ListReviewsResponseDTO struct {
	Items []domainreview.Review `json:"items"`
	Total int64                 `json:"total"`
}

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
// @Summary Submeter avaliação
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body review.SubmitReviewInput true "Dados da avaliação"
// @Success 201 {object} domainreview.Review "Avaliação criada"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 403 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /reviews [post]
func (h *ReviewHandler) Submit(c *gin.Context) {
	ownerID, problems := extractUserIDProblems(c)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ut := extractUserType(c); ut != "" && ut != domainuser.UserTypeOwner {
		h.errorService.RespondWithProblems(c, []exceptions.ProblemDetails{{
			Type:   string(exceptions.Forbidden),
			Title:  "Apenas donos podem registrar avaliações",
			Status: http.StatusForbidden,
		}}, "forbidden", http.StatusForbidden)
		return
	}
	var req submitReviewRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}
	reqID, err := uuid.Parse(req.RequestID)
	if err != nil {
		h.errorService.RespondWithProblems(c, []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "ID da solicitação inválido",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}, "invalid_request_id", http.StatusBadRequest)
		return
	}

	reviewCreated, problems := h.uc.Submit.Execute(c.Request.Context(), review.SubmitReviewInput{
		RequestID: reqID,
		OwnerID:   ownerID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "submit_review_failed", http.StatusBadRequest)
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
// @Summary Listar avaliações do prestador
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param provider_id path string true "ID do prestador"
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} http.ListReviewsResponseDTO "Lista de avaliações"
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 403 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 404 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 409 {object} exceptions.ProblemDetailsOutputDTO
// @Failure 500 {object} exceptions.ProblemDetailsOutputDTO
// @Router /reviews/provider/{provider_id} [get]
func (h *ReviewHandler) ListForProvider(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}

	// Se o usuário autenticado for prestador, ele só pode listar o próprio perfil.
	if extractUserType(c) == domainuser.UserTypeProvider {
		authProviderID, problems := providerIDFromContextProblems(c, h.providerRepo, true)
		if len(problems) > 0 {
			h.errorService.RespondWithProblems(c, problems, "forbidden", http.StatusForbidden)
			return
		}
		if authProviderID != providerID {
			h.errorService.RespondWithProblems(c, []exceptions.ProblemDetails{{
				Type:   string(exceptions.Forbidden),
				Title:  "Não autorizado a listar avaliações de outro prestador",
				Status: http.StatusForbidden,
			}}, "forbidden", http.StatusForbidden)
			return
		}
	}

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	reviews, total, problems := h.uc.ListForProvider.Execute(c.Request.Context(), review.ListReviewsForProviderInput{ProviderID: providerID, Page: page, Limit: limit})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "list_reviews_failed", http.StatusBadRequest)
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
