package usecases

import (
	"context"
	"errors"
	"math"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListReviewsInput struct {
	ProviderID string `json:"provider_id"`
	UserID     string `json:"user_id"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type ListReviewsOutput struct {
	Reviews      []*entities.Review `json:"reviews"`
	TotalItems   int64              `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	ItemsPerPage int                `json:"items_per_page"`
}

type ListReviewsUseCase struct {
	reviewRepository entities.ReviewRepository
	logger           logging.LoggerInterface
}

func NewListReviewsUseCase(reviewRepository entities.ReviewRepository, logger logging.LoggerInterface) *ListReviewsUseCase {
	return &ListReviewsUseCase{
		reviewRepository: reviewRepository,
		logger:           logger,
	}
}

func (uc *ListReviewsUseCase) Execute(ctx context.Context, input ListReviewsInput) (*ListReviewsOutput, []exceptions.ProblemDetails) {
	const from = "ListReviewsUseCase.Execute"

	if input.ProviderID == "" && input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Filtro ausente", errors.New("Informe provider_id ou user_id"))
	}

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	reviews, total, err := uc.reviewRepository.List(input.ProviderID, input.UserID, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar reviews", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	return &ListReviewsOutput{
		Reviews:      reviews,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}, nil
}
