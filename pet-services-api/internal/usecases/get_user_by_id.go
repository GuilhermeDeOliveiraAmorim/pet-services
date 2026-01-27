package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type GetUserByIDInput struct {
	UserID string `json:"user_id"`
}

type GetUserByIDOutput struct {
	User *entities.User `json:"user"`
}

type GetUserByIDUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewGetUserByIDUseCase(userRepository entities.UserRepository, logger logging.LoggerInterface) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (uc *GetUserByIDUseCase) Execute(ctx context.Context, input GetUserByIDInput) (*GetUserByIDOutput, []exceptions.ProblemDetails) {
	const from = "GetUserByIDUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório para buscar"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	return &GetUserByIDOutput{
		User: user,
	}, nil
}
