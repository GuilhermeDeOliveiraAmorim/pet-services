package review

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainReview "github.com/guilherme/pet-services-api/internal/domain/review"
)

// ListReviewsForProviderUseCase lista avaliações de um prestador.
type ListReviewsForProviderUseCase struct {
	reviewRepo domainReview.Repository
	logger     *slog.Logger
}

func NewListReviewsForProviderUseCase(reviewRepo domainReview.Repository, logger *slog.Logger) *ListReviewsForProviderUseCase {
	return &ListReviewsForProviderUseCase{reviewRepo: reviewRepo, logger: logging.EnsureLogger(logger)}
}

// ListReviewsForProviderInput entrada para listagem.
type ListReviewsForProviderInput struct {
	ProviderID uuid.UUID
	Page       int
	Limit      int
}

func (uc *ListReviewsForProviderUseCase) Execute(ctx context.Context, input ListReviewsForProviderInput) ([]*domainReview.Review, int64, error) {
	var (
		reviews []*domainReview.Review
		total   int64
		err     error
	)
	defer logging.UseCase(ctx, uc.logger, "ListReviewsForProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	if input.ProviderID == uuid.Nil {
		err = fmt.Errorf("providerID é obrigatório")
		return nil, 0, err
	}

	page, limit := normalizePagination(input.Page, input.Limit)

	reviews, total, err = uc.reviewRepo.ListByProvider(ctx, input.ProviderID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// normalizePagination aplica defaults e evita valores inválidos.
func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
