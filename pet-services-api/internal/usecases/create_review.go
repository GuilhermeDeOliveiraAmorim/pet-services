package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CreateReviewInput struct {
	UserID     string  `json:"user_id"`
	ProviderID string  `json:"provider_id"`
	Rating     float64 `json:"rating"`
	Comment    string  `json:"comment"`
}

type CreateReviewInputBody struct {
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}

type CreateReviewOutput struct {
	Message string          `json:"message,omitempty"`
	Detail  string          `json:"detail,omitempty"`
	Review  *entities.Review `json:"review,omitempty"`
}

type CreateReviewUseCase struct {
	userRepository    entities.UserRepository
	providerRepository entities.ProviderRepository
	requestRepository entities.RequestRepository
	reviewRepository  entities.ReviewRepository
	logger            logging.LoggerInterface
}

func NewCreateReviewUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	reviewRepository entities.ReviewRepository,
	logger logging.LoggerInterface,
) *CreateReviewUseCase {
	return &CreateReviewUseCase{
		userRepository:    userRepository,
		providerRepository: providerRepository,
		requestRepository: requestRepository,
		reviewRepository:  reviewRepository,
		logger:            logger,
	}
}

func (uc *CreateReviewUseCase) Execute(ctx context.Context, input CreateReviewInput) (*CreateReviewOutput, []exceptions.ProblemDetails) {
	const from = "CreateReviewUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ProviderID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do provedor ausente", errors.New("O ID do provedor é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem criar reviews"))
	}

	_, err = uc.providerRepository.FindByID(input.ProviderID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	exists, err := uc.requestRepository.ExistsCompleted(user.ID, input.ProviderID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar requests concluídas", err)
	}
	if !exists {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("É necessário ter uma request concluída com o provedor"))
	}

	review, problems := entities.NewReview(user.ID, input.ProviderID, input.Rating, input.Comment)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Review inválida", problems)
		return nil, problems
	}

	if err := uc.reviewRepository.Create(review); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar review", err)
	}

	return &CreateReviewOutput{
		Message: "Review criada com sucesso",
		Detail:  "A review foi registrada",
		Review:  review,
	}, nil
}
