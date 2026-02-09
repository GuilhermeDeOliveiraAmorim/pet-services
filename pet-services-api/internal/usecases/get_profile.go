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
	User *entities.User `json:"user"`
}

type GetProfileUseCase struct {
	userRepository entities.UserRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewGetProfileUseCase(userRepo entities.UserRepository, storageService storage.ObjectStorage, logger logging.LoggerInterface) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepository: userRepo,
		storage:        storageService,
		logger:         logger,
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

	if err := signUserPhotos(ctx, uc.storage, user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URLs das fotos", err)
	}

	return &GetProfileOutput{
		User: user,
	}, nil
}
