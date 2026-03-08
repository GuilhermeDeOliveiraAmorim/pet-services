package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type GetProfileInput struct {
	UserID string `json:"user_id"`
}

type GetProfileOutput struct {
	User       *entities.User `json:"user"`
	ProviderID string         `json:"provider_id,omitempty"`
}

type GetProfileUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewGetProfileUseCase(
	userRepo entities.UserRepository,
	providerRepo entities.ProviderRepository,
	storageService storage.ObjectStorage,
	logger logging.LoggerInterface,
) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepository:     userRepo,
		providerRepository: providerRepo,
		storage:            storageService,
		logger:             logger,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (*GetProfileOutput, []exceptions.ProblemDetails) {
	const from = "GetProfileUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := storage.SignUserPhotos(ctx, uc.storage, user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URLs das fotos", err)
	}

	output := &GetProfileOutput{User: user}

	if user.IsProvider() {
		provider, err := uc.providerRepository.FindByUserID(user.ID)
		if err != nil {
			if err.Error() != consts.ProviderNotFoundError {
				return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor do usuário", err)
			}
		} else if provider != nil {
			output.ProviderID = provider.ID
		}
	}

	return output, nil
}
