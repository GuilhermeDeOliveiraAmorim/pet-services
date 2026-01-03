package review

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainReview "github.com/guilherme/pet-services-api/internal/domain/review"
)

// ListReviewsForProviderUseCase lista avaliações de um prestador.
type ListReviewsForProviderUseCase struct {
	reviewRepo domainReview.Repository
	logger     logging.LoggerService
}

func NewListReviewsForProviderUseCase(reviewRepo domainReview.Repository, logger logging.LoggerService) *ListReviewsForProviderUseCase {
	return &ListReviewsForProviderUseCase{reviewRepo: reviewRepo, logger: logger}
}

// ListReviewsForProviderInput entrada para listagem.
type ListReviewsForProviderInput struct {
	ProviderID uuid.UUID
	Page       int
	Limit      int
}

const LIST_REVIEWS_USECASE = "LIST_REVIEWS_USECASE"

func (uc *ListReviewsForProviderUseCase) Execute(ctx context.Context, input ListReviewsForProviderInput) ([]*domainReview.Review, int64, []exceptions.ProblemDetails) {
       uc.logger.Log(logging.Logger{
	       Context: ctx,
	       TypeLog: logging.LoggerTypes.INFO,
	       Layer:   logging.LoggerLayers.USECASES,
	       Code:    exceptions.RFC200_CODE,
	       From:    LIST_REVIEWS_USECASE,
	       Message: logging.DEFAULTMESSAGES.START,
       })

       if input.ProviderID == uuid.Nil {
	       uc.logger.Log(logging.Logger{
		       Context: ctx,
		       TypeLog: logging.LoggerTypes.ERROR,
		       Layer:   logging.LoggerLayers.USECASES,
		       Code:    exceptions.RFC400_CODE,
		       From:    LIST_REVIEWS_USECASE,
		       Message: "ProviderID é obrigatório",
		       Error:   errors.New("providerID é obrigatório"),
	       })
	       return nil, 0, []exceptions.ProblemDetails{{
		       Type:   exceptions.RFC400,
		       Title:  "ProviderID é obrigatório",
		       Status: exceptions.RFC400_CODE,
		       Detail: "O ID do prestador é obrigatório.",
	       }}
       }

       page, limit := normalizePagination(input.Page, input.Limit)

       reviews, total, err := uc.reviewRepo.ListByProvider(ctx, input.ProviderID, page, limit)
       if err != nil {
	       uc.logger.Log(logging.Logger{
		       Context: ctx,
		       TypeLog: logging.LoggerTypes.ERROR,
		       Layer:   logging.LoggerLayers.USECASES,
		       Code:    exceptions.RFC500_CODE,
		       From:    LIST_REVIEWS_USECASE,
		       Message: "Erro ao listar avaliações",
		       Error:   err,
	       })
	       return nil, 0, []exceptions.ProblemDetails{{
		       Type:   exceptions.RFC500,
		       Title:  "Erro ao listar avaliações",
		       Status: exceptions.RFC500_CODE,
		       Detail: err.Error(),
	       }}
       }

       uc.logger.Log(logging.Logger{
	       Context: ctx,
	       TypeLog: logging.LoggerTypes.INFO,
	       Layer:   logging.LoggerLayers.USECASES,
	       Code:    exceptions.RFC200_CODE,
	       From:    LIST_REVIEWS_USECASE,
	       Message: logging.DEFAULTMESSAGES.END,
       })

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
