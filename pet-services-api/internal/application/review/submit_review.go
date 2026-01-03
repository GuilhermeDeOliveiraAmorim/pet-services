package review

import (
    "context"
    "fmt"
    "strings"

    "github.com/google/uuid"

    domainProvider "github.com/guilherme/pet-services-api/internal/domain/provider"
    domainRequest "github.com/guilherme/pet-services-api/internal/domain/request"
    domainReview "github.com/guilherme/pet-services-api/internal/domain/review"
)

// SubmitReviewUseCase cria uma avaliação para uma solicitação concluída.
type SubmitReviewUseCase struct {
    reviewRepo   domainReview.Repository
    requestRepo  domainRequest.Repository
    providerRepo domainProvider.Repository
}

func NewSubmitReviewUseCase(reviewRepo domainReview.Repository, requestRepo domainRequest.Repository, providerRepo domainProvider.Repository) *SubmitReviewUseCase {
    return &SubmitReviewUseCase{
        reviewRepo:   reviewRepo,
        requestRepo:  requestRepo,
        providerRepo: providerRepo,
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
    if input.RequestID == uuid.Nil {
        return nil, fmt.Errorf("requestID é obrigatório")
    }
    if input.OwnerID == uuid.Nil {
        return nil, fmt.Errorf("ownerID é obrigatório")
    }

    req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
    if err != nil {
        return nil, err
    }

    if req.OwnerID != input.OwnerID {
        return nil, fmt.Errorf("não autorizado a avaliar esta solicitação")
    }
    if req.Status != domainRequest.StatusCompleted {
        return nil, domainReview.ErrRequestNotCompleted
    }

    exists, err := uc.reviewRepo.ExistsByRequestID(ctx, input.RequestID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, domainReview.ErrReviewAlreadyExists
    }

    rating := input.Rating
    if rating < 1 || rating > 5 {
        return nil, domainReview.ErrInvalidRating
    }

    comment := strings.TrimSpace(input.Comment)
    review, err := domainReview.NewReview(req.ID, req.ProviderID, req.OwnerID, rating, comment)
    if err != nil {
        return nil, err
    }

    if err := uc.reviewRepo.Create(ctx, review); err != nil {
        return nil, fmt.Errorf("falha ao salvar avaliação: %w", err)
    }

    // Atualiza métricas do prestador.
    provider, err := uc.providerRepo.FindByID(ctx, req.ProviderID)
    if err == nil { // se não encontrar, apenas ignora atualização de métricas
        provider.UpdateRating(float64(rating))
        _ = uc.providerRepo.Update(ctx, provider)
    }

    return review, nil
}
