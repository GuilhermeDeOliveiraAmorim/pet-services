package review

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainProvider "pet-services-api/internal/domain/provider"
	domainRequest "pet-services-api/internal/domain/request"
	domainReview "pet-services-api/internal/domain/review"
)

// SubmitReviewUseCase cria uma avaliação para uma solicitação concluída.
type SubmitReviewUseCase struct {
	reviewRepo   domainReview.Repository
	requestRepo  domainRequest.Repository
	providerRepo domainProvider.Repository
	logger       logging.LoggerService
}

func NewSubmitReviewUseCase(reviewRepo domainReview.Repository, requestRepo domainRequest.Repository, providerRepo domainProvider.Repository, logger logging.LoggerService) *SubmitReviewUseCase {
	return &SubmitReviewUseCase{
		reviewRepo:   reviewRepo,
		requestRepo:  requestRepo,
		providerRepo: providerRepo,
		logger:       logger,
	}
}

// SubmitReviewInput dados para registrar avaliação.
type SubmitReviewInput struct {
	RequestID uuid.UUID
	OwnerID   uuid.UUID
	Rating    int
	Comment   string
}

const SUBMIT_REVIEW_USECASE = "SUBMIT_REVIEW_USECASE"

func (uc *SubmitReviewUseCase) Execute(ctx context.Context, input SubmitReviewInput) (*domainReview.Review, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    SUBMIT_REVIEW_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.RequestID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "RequestID é obrigatório",
			Error:   errors.New("requestID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "RequestID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID da solicitação é obrigatório.",
		}}
	}
	if input.OwnerID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "OwnerID é obrigatório",
			Error:   errors.New("ownerID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "OwnerID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do dono é obrigatório.",
		}}
	}

	req, err := uc.requestRepo.FindByID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Solicitação não encontrada",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Solicitação não encontrada",
			Status: exceptions.RFC404_CODE,
			Detail: "A solicitação informada não foi encontrada.",
		}}
	}

	if req.OwnerID != input.OwnerID {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC403_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Não autorizado a avaliar esta solicitação",
			Error:   errors.New("não autorizado a avaliar esta solicitação"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC403,
			Title:  "Não autorizado",
			Status: exceptions.RFC403_CODE,
			Detail: "Você não tem permissão para avaliar esta solicitação.",
		}}
	}
	if req.Status != domainRequest.StatusCompleted {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Solicitação não está concluída",
			Error:   domainReview.ErrRequestNotCompleted,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Solicitação não está concluída",
			Status: exceptions.RFC409_CODE,
			Detail: "A solicitação precisa estar concluída para ser avaliada.",
		}}
	}

	exists, err := uc.reviewRepo.ExistsByRequestID(ctx, input.RequestID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Erro ao verificar avaliação existente",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao verificar avaliação existente",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}
	if exists {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Avaliação já existe para esta solicitação",
			Error:   domainReview.ErrReviewAlreadyExists,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Avaliação já existe",
			Status: exceptions.RFC409_CODE,
			Detail: "Já existe uma avaliação para esta solicitação.",
		}}
	}

	rating := input.Rating
	if rating < 1 || rating > 5 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Nota inválida",
			Error:   domainReview.ErrInvalidRating,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Nota inválida",
			Status: exceptions.RFC400_CODE,
			Detail: "A nota deve ser entre 1 e 5.",
		}}
	}

	comment := strings.TrimSpace(input.Comment)
	result, err := domainReview.NewReview(req.ID, req.ProviderID, req.OwnerID, rating, comment)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Dados inválidos para avaliação",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Dados inválidos para avaliação",
			Status: exceptions.RFC400_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.reviewRepo.Create(ctx, result); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    SUBMIT_REVIEW_USECASE,
			Message: "Falha ao salvar avaliação",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar avaliação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	// Atualiza métricas do prestador.
	provider, err := uc.providerRepo.FindByID(ctx, req.ProviderID)
	if err == nil { // se não encontrar, apenas ignora atualização de métricas
		provider.UpdateRating(float64(rating))
		_ = uc.providerRepo.Update(ctx, provider)
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    SUBMIT_REVIEW_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return result, nil
}
