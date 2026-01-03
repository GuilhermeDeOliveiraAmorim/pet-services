package review

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
	domainReview "github.com/guilherme/pet-services-api/internal/domain/review"
)

// SubmitReviewUseCase cria uma avaliação para uma solicitação concluída.
type SubmitReviewUseCase struct {
	reviewRepo   domainReview.Repository
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       *slog.Logger
}

func NewSubmitReviewUseCase(reviewRepo domainReview.Repository, requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger *slog.Logger) *SubmitReviewUseCase {
	return &SubmitReviewUseCase{
		reviewRepo:   reviewRepo,
		requestRepo:  requestRepo,
		providerRepo: providerRepo,
		logger:       logging.EnsureLogger(logger),
	}
}

// SubmitReviewInput dados para registrar avaliação.
type SubmitReviewInput struct {
	RequestID uuid.UUID
	OwnerID   uuid.UUID
	Rating    int
	Comment   string
}

func (uc *SubmitReviewUseCase) Execute(ctx context.Context, input SubmitReviewInput) (*domainReview.Review, error) {
	var (
		result *domainReview.Review
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "SubmitReviewUseCase", slog.String("request_id", input.RequestID.String()), slog.String("owner_id", input.OwnerID.String()))(&err)

	if input.RequestID == uuid.Nil {
		err = fmt.Errorf("requestID é obrigatório")
		return nil, err
	}
	if input.OwnerID == uuid.Nil {
		err = fmt.Errorf("ownerID é obrigatório")
		return nil, err
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	if req.OwnerID != input.OwnerID {
		err = fmt.Errorf("não autorizado a avaliar esta solicitação")
		return nil, err
	}
	if req.Status != domainRequest.StatusCompleted {
		err = domainReview.ErrRequestNotCompleted
		return nil, err
	}

	exists, err := uc.reviewRepo.ExistsByRequestID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		err = domainReview.ErrReviewAlreadyExists
		return nil, err
	}

	rating := input.Rating
	if rating < 1 || rating > 5 {
		err = domainReview.ErrInvalidRating
		return nil, err
	}

	comment := strings.TrimSpace(input.Comment)
	result, err = domainReview.NewReview(req.ID, req.ProviderID, req.OwnerID, rating, comment)
	if err != nil {
		return nil, err
	}

	if err := uc.reviewRepo.Create(ctx, result); err != nil {
		err = fmt.Errorf("falha ao salvar avaliação: %w", err)
		return nil, err
	}

	// Atualiza métricas do prestador.
	provider, err := uc.providerRepo.FindByID(ctx, req.ProviderID)
	if err == nil { // se não encontrar, apenas ignora atualização de métricas
		provider.UpdateRating(float64(rating))
		_ = uc.providerRepo.Update(ctx, provider)
	}

	return result, nil
}
