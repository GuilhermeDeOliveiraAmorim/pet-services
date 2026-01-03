package review

import (
    "context"
    "fmt"

    "github.com/google/uuid"

    domainReview "github.com/guilherme/pet-services-api/internal/domain/review"
)

// ListReviewsForProviderUseCase lista avaliações de um prestador.
type ListReviewsForProviderUseCase struct {
    reviewRepo domainReview.Repository
}

func NewListReviewsForProviderUseCase(reviewRepo domainReview.Repository) *ListReviewsForProviderUseCase {
    return &ListReviewsForProviderUseCase{reviewRepo: reviewRepo}
}

// ListReviewsForProviderInput entrada para listagem.
type ListReviewsForProviderInput struct {
    ProviderID uuid.UUID
    Page       int
    Limit      int
}

func (uc *ListReviewsForProviderUseCase) Execute(ctx context.Context, input ListReviewsForProviderInput) ([]*domainReview.Review, int64, error) {
    if input.ProviderID == uuid.Nil {
        return nil, 0, fmt.Errorf("providerID é obrigatório")
    }

    page, limit := normalizePagination(input.Page, input.Limit)

    reviews, total, err := uc.reviewRepo.ListByProvider(ctx, input.ProviderID, page, limit)
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
